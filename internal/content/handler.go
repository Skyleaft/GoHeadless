package content

import (
	"context"

	"GoHeadless/internal/domain"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	group := router.Group("/content/:collSlug")
	group.Get("/", h.ListRecords)
	group.Post("/", h.CreateRecord)
	group.Get("/:id", h.GetRecord)
	group.Put("/:id", h.UpdateRecord)
	group.Delete("/:id", h.DeleteRecord)
}

func (h *Handler) ListRecords(c fiber.Ctx) error {
	collSlug := c.Params("collSlug")
	ctx := context.Background()
	records, err := h.service.GetRecords(ctx, collSlug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

func (h *Handler) CreateRecord(c fiber.Ctx) error {
	collSlug := c.Params("collSlug")
	var record domain.Record
	if err := c.Bind().Body(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	id, err := h.service.CreateRecord(ctx, collSlug, record)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (h *Handler) GetRecord(c fiber.Ctx) error {
	collSlug := c.Params("collSlug")
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid object id"})
	}

	ctx := context.Background()
	record, err := h.service.GetRecord(ctx, collSlug, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "record not found"})
	}
	return c.JSON(record)
}

func (h *Handler) UpdateRecord(c fiber.Ctx) error {
	collSlug := c.Params("collSlug")
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid object id"})
	}

	var record domain.Record
	if err := c.Bind().Body(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	if err := h.service.UpdateRecord(ctx, collSlug, id, record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DeleteRecord(c fiber.Ctx) error {
	collSlug := c.Params("collSlug")
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid object id"})
	}

	ctx := context.Background()
	if err := h.service.DeleteRecord(ctx, collSlug, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
