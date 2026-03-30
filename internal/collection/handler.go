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

func (h *Handler) ListCollections(c fiber.Ctx) error {
	ctx := context.Background()
	colls, err := h.service.GetAllCollections(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(colls)
}

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

func (h *Handler) GetCollection(c fiber.Ctx) error {
	slug := c.Params("slug")
	ctx := context.Background()
	coll, err := h.service.GetCollection(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "collection not found"})
	}
	return c.JSON(coll)
}

func (h *Handler) DeleteCollection(c fiber.Ctx) error {
	slug := c.Params("slug")
	ctx := context.Background()
	if err := h.service.DeleteCollection(ctx, slug); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
