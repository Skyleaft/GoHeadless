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

// AuthorizeCollection checks if the authenticated user has access to the requested collection
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

	// 2. Superuser Bypass
	isSuperuser, _ := c.Locals("is_superuser").(bool)
	if isSuperuser {
		return c.Next()
	}

	// 3. Check Public Access (Read only)
	method := c.Method()
	if coll.Access != nil && coll.Access.IsPublic && method == fiber.MethodGet {
		return c.Next()
	}

	// 4. Require Auth for everything else
	roleID, _ := c.Locals("role_id").(string)
	if roleID == "" {
		// If we reached here without a role_id, it means Authenticate middleware was skipped or failed
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authentication required"})
	}

	// 5. Check Granular Permissions (CRUD Policy)
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
