package collection

import (
	"context"

	"GoHeadless/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes(router fiber.Router) {
	group := router.Group("/collections")
	group.Get("/", h.ListCollections)
	group.Post("/", h.CreateCollection)
	group.Get("/:slug", h.GetCollection)
	group.Delete("/:slug", h.DeleteCollection)
}

// ListCollections list all collections
// @Summary List all collections
// @Description Get metadata for all user-defined collections
// @Tags collections
// @Produce json
// @Success 200 {array} domain.Collection
// @Router /collections [get]
func (h *Handler) ListCollections(c fiber.Ctx) error {
	ctx := context.Background()
	colls, err := h.service.GetAllCollections(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(colls)
}

// CreateCollection create a new collection
// @Summary Create a new collection
// @Description Define a new dynamic collection with fields
// @Tags collections
// @Accept json
// @Produce json
// @Param collection body domain.Collection true "Collection Definition"
// @Success 201 {object} domain.Collection
// @Router /collections [post]
func (h *Handler) CreateCollection(c fiber.Ctx) error {
	var coll domain.Collection
	if err := c.Bind().Body(&coll); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	if err := h.service.CreateCollection(ctx, &coll); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(coll)
}

// GetCollection get a collection by slug
// @Summary Get collection by slug
// @Description Get metadata for a specific collection by its unique slug
// @Tags collections
// @Produce json
// @Param slug path string true "Collection Slug"
// @Success 200 {object} domain.Collection
// @Router /collections/{slug} [get]
func (h *Handler) GetCollection(c fiber.Ctx) error {
	slug := c.Params("slug")
	ctx := context.Background()
	coll, err := h.service.GetCollection(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "collection not found"})
	}
	return c.JSON(coll)
}

// DeleteCollection delete a collection by slug
// @Summary Delete collection by slug
// @Description Remove a collection definition by its slug
// @Tags collections
// @Param slug path string true "Collection Slug"
// @Success 204 "No Content"
// @Router /collections/{slug} [delete]
func (h *Handler) DeleteCollection(c fiber.Ctx) error {
	slug := c.Params("slug")
	ctx := context.Background()
	if err := h.service.DeleteCollection(ctx, slug); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
