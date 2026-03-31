package setup

import (
	"github.com/gofiber/fiber/v3"
)

// SetupRequest is the payload for the initial superadmin registration
type SetupRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes(router fiber.Router) {
	router.Get("/setup/status", h.Status)
	router.Post("/setup", h.Initialize)
}

// Status godoc
// @Summary Check if system setup is required
// @Tags setup
// @Produce json
// @Success 200 {object} map[string]bool
// @Router /setup/status [get]
func (h *Handler) Status(c fiber.Ctx) error {
	required, err := h.svc.IsSetupRequired(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"setup_required": required})
}

// Initialize godoc
// @Summary Create the initial superadmin user
// @Description This endpoint only works when no users exist in the system. Calling it after initial setup returns a 403.
// @Tags setup
// @Accept json
// @Produce json
// @Param body body SetupRequest true "Initial credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /setup [post]
func (h *Handler) Initialize(c fiber.Ctx) error {
	// Guard: only allow if setup is still required
	required, err := h.svc.IsSetupRequired(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !required {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "system is already initialized"})
	}

	var req SetupRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := h.svc.Bootstrap(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "System initialized successfully",
		"username": user.Username,
	})
}
