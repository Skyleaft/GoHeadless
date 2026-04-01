package middleware

import (
	"errors"
	"strings"

	"GoHeadless/internal/apierr"
	"GoHeadless/internal/auth"
	"GoHeadless/internal/collection"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

// CollectionUnknownLocalsKey is set when GET /content/:slug (list) is allowed without metadata (empty list).
const CollectionUnknownLocalsKey = "collection_unknown"

type RBACMiddleware struct {
	authSvc auth.Service
	collSvc collection.Service
}

func NewRBACMiddleware(authSvc auth.Service, collSvc collection.Service) *RBACMiddleware {
	return &RBACMiddleware{
		authSvc: authSvc,
		collSvc: collSvc,
	}
}

// Authenticate extracts and validates the JWT token from the Authorization header or cookie
func (m *RBACMiddleware) Authenticate(c fiber.Ctx) error {
	// 1. Get token from Header or Cookie
	authHeader := c.Get("Authorization")
	tokenStr := ""

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		tokenStr = c.Cookies("jwt")
	}

	if tokenStr == "" {
		return apierr.Unauthorized("missing authentication token")
	}

	// 2. Validate token
	claims, err := m.authSvc.ValidateToken(tokenStr)
	if err != nil {
		return apierr.Unauthorized(err.Error())
	}

	// 3. Store claims in context
	c.Locals("user_id", claims.UserID)
	c.Locals("role_id", claims.RoleID)
	c.Locals("is_superuser", claims.IsSuperuser)

	return c.Next()
}

// RequireSuperadmin ensures the authenticated user is a superadmin
func (m *RBACMiddleware) RequireSuperadmin(c fiber.Ctx) error {
	isSuperuser, _ := c.Locals("is_superuser").(bool)
	if !isSuperuser {
		return apierr.Forbidden("superuser privileges required")
	}
	return c.Next()
}

// tryAuthenticate sets user locals when a valid JWT is present; never fails the request.
func (m *RBACMiddleware) tryAuthenticate(c fiber.Ctx) {
	authHeader := c.Get("Authorization")
	tokenStr := ""
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		tokenStr = c.Cookies("jwt")
	}
	if tokenStr == "" {
		return
	}
	claims, err := m.authSvc.ValidateToken(tokenStr)
	if err != nil {
		return
	}
	c.Locals("user_id", claims.UserID)
	c.Locals("role_id", claims.RoleID)
	c.Locals("is_superuser", claims.IsSuperuser)
}

// AuthorizeCollection checks access for /collections/... routes (metadata must exist).
func (m *RBACMiddleware) AuthorizeCollection(c fiber.Ctx) error {
	return m.authorizeCollection(c, false)
}

// AuthorizeContentCollection is for /content/:slug/... only. Unknown slug: GET list returns empty data;
// POST create requires auth and inserts without a schema document (see CreateRecord service).
func (m *RBACMiddleware) AuthorizeContentCollection(c fiber.Ctx) error {
	return m.authorizeCollection(c, true)
}

func (m *RBACMiddleware) authorizeCollection(c fiber.Ctx, allowContentWithoutMetadata bool) error {
	slug := c.Params("slug")
	if slug == "" {
		slug = c.Params("collSlug")
	}
	if slug == "" {
		return c.Next()
	}

	coll, err := m.collSvc.GetCollection(c.Context(), slug)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) && allowContentWithoutMetadata {
			if isContentListGET(c) {
				c.Locals(CollectionUnknownLocalsKey, true)
				return c.Next()
			}
			// Any write or single-record op on unknown collection — require auth then proceed
			if err := m.Authenticate(c); err != nil {
				return err
			}
			return c.Next()
		}
		return apierr.NotFound("collection not found")
	}

	method := c.Method()

	if coll.Access != nil && coll.Access.IsPublic && method == fiber.MethodGet {
		m.tryAuthenticate(c)
		return c.Next()
	}

	hasAuth := c.Locals("user_id") != nil
	if !hasAuth {
		if err := m.Authenticate(c); err != nil {
			return err
		}
	}

	isSuperuser, _ := c.Locals("is_superuser").(bool)
	if isSuperuser {
		return c.Next()
	}

	// No access policy defined — allow any authenticated user
	if coll.Access == nil {
		return c.Next()
	}

	roleID, _ := c.Locals("role_id").(string)
	var allowedRoles []string
	switch method {
	case fiber.MethodPost:
		allowedRoles = coll.Access.CRUDPolicy.Create
	case fiber.MethodGet:
		allowedRoles = coll.Access.CRUDPolicy.Read
	case fiber.MethodPut, fiber.MethodPatch:
		allowedRoles = coll.Access.CRUDPolicy.Update
	case fiber.MethodDelete:
		allowedRoles = coll.Access.CRUDPolicy.Delete
	}

	for _, r := range allowedRoles {
		if r == roleID {
			return c.Next()
		}
	}

	return apierr.Forbidden("insufficient permissions for this collection")
}

// isContentListGET is true for GET /content/:slug (list), false for GET /content/:slug/:id.
func isContentListGET(c fiber.Ctx) bool {
	return c.Method() == fiber.MethodGet && c.Params("id") == ""
}
