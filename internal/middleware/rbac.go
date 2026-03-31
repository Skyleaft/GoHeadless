package middleware

import (
	"strings"

	"GoHeadless/internal/auth"
	"GoHeadless/internal/collection"

	"github.com/gofiber/fiber/v3"
)

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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authentication token"})
	}

	// 2. Validate token
	claims, err := m.authSvc.ValidateToken(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
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
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "superuser privileges required"})
	}
	return c.Next()
}

// AuthorizeCollection checks if the authenticated user has access to the requested collection.
// For public collections with GET requests, it allows access without authentication.
// For all other cases, it requires and validates authentication.
func (m *RBACMiddleware) AuthorizeCollection(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		slug = c.Params("collSlug")
	}
	if slug == "" {
		return c.Next()
	}

	// 1. Get Collection Permissions
	coll, err := m.collSvc.GetCollection(c.Context(), slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "collection not found"})
	}

	method := c.Method()

	// 2. Check Public Access (Read only) - allows bypass of authentication
	// If Access is nil or IsPublic is true, allow GET requests without authentication
	if coll.Access != nil && coll.Access.IsPublic && method == fiber.MethodGet {
		return c.Next()
	}

	// 3. For non-public collections or non-GET methods, require authentication
	// Only process auth if it hasn't been done already (check for user_id in locals)
	hasAuth := c.Locals("user_id") != nil
	if !hasAuth {
		// Run authentication inline
		if err := m.Authenticate(c); err != nil {
			return err
		}
	}

	// 4. Superuser Bypass (check again after auth)
	isSuperuser, _ := c.Locals("is_superuser").(bool)
	if isSuperuser {
		return c.Next()
	}

	// 5. Check Granular Permissions (CRUD Policy)
	roleID, _ := c.Locals("role_id").(string)
	if coll.Access != nil {
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
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions for this collection"})
}
