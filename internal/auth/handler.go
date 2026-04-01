package auth

import (
	"GoHeadless/internal/apierr"
	"GoHeadless/internal/domain"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) PublicRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/login", h.Login)
}

func (h *Handler) AdminRoutes(router fiber.Router) {
	router.Get("/users", h.GetUsers)
	router.Post("/users", h.CreateUser)
	router.Delete("/users/:id", h.DeleteUser)

	router.Get("/roles", h.GetRoles)
	router.Post("/roles", h.CreateRole)
	router.Put("/roles/:id", h.UpdateRole)
	router.Delete("/roles/:id", h.DeleteRole)

	router.Get("/stats", h.GetStats)
}

// Login godoc
// @Summary Login user
// @Description Authenticate with username and password, returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apierr.BadRequest("invalid request body")
	}

	resp, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return apierr.Unauthorized(err.Error())
	}
	return c.JSON(resp)
}

func (h *Handler) GetUsers(c fiber.Ctx) error {
	users, err := h.svc.GetAllUsers(c.Context())
	if err != nil {
		return apierr.Internal("failed to get users", err)
	}
	return c.JSON(users)
}

func (h *Handler) CreateUser(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apierr.BadRequest("invalid request body")
	}

	user, err := h.svc.Register(c.Context(), req)
	if err != nil {
		return apierr.BadRequest(err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *Handler) DeleteUser(c fiber.Ctx) error {
	if err := h.svc.DeleteUser(c.Context(), c.Params("id")); err != nil {
		return apierr.Internal("failed to delete user", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetRoles(c fiber.Ctx) error {
	roles, err := h.svc.GetAllRoles(c.Context())
	if err != nil {
		return apierr.Internal("failed to get roles", err)
	}
	return c.JSON(roles)
}

func (h *Handler) CreateRole(c fiber.Ctx) error {
	var role domain.Role
	if err := c.Bind().JSON(&role); err != nil {
		return apierr.BadRequest("invalid request body")
	}

	if err := h.svc.CreateRole(c.Context(), &role); err != nil {
		return apierr.BadRequest(err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(role)
}

func (h *Handler) UpdateRole(c fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return apierr.BadRequest("invalid role id")
	}

	var role domain.Role
	if err := c.Bind().JSON(&role); err != nil {
		return apierr.BadRequest("invalid request body")
	}
	role.ID = oid

	if err := h.svc.UpdateRole(c.Context(), &role); err != nil {
		return apierr.Internal("failed to update role", err)
	}
	return c.JSON(role)
}

func (h *Handler) DeleteRole(c fiber.Ctx) error {
	if err := h.svc.DeleteRole(c.Context(), c.Params("id")); err != nil {
		return apierr.Internal("failed to delete role", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetStats(c fiber.Ctx) error {
	stats, err := h.svc.GetStats(c.Context())
	if err != nil {
		return apierr.Internal("failed to get stats", err)
	}
	return c.JSON(stats)
}
