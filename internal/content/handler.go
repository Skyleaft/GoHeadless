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
	router.Get("/", h.ListRecords)
	router.Post("/", h.CreateRecord)
	router.Get("/:id", h.GetRecord)
	router.Put("/:id", h.UpdateRecord)
	router.Delete("/:id", h.DeleteRecord)
}

// ListRecords list all records for a collection
// @Summary List all records in a collection
// @Description Fetch all dynamic content records for the given collection slug
// @Tags content
// @Produce json
// @Param collSlug path string true "Collection Slug"
// @Success 200 {array} map[string]interface{}
// @Router /content/{collSlug} [get]
func (h *Handler) ListRecords(c fiber.Ctx) error {
	slug := c.Params("slug")
	ctx := context.Background()
	records, err := h.service.GetRecords(ctx, slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(records)
}

// CreateRecord create a new record in a collection
// @Summary Create a new record
// @Description Upsert dynamic content into a collection. Validates against collection schema.
// @Tags content
// @Accept json
// @Produce json
// @Param collSlug path string true "Collection Slug"
// @Param record body map[string]interface{} true "Dynamic Record Data"
// @Success 201 {object} map[string]interface{}
// @Router /content/{collSlug} [post]
func (h *Handler) CreateRecord(c fiber.Ctx) error {
	slug := c.Params("slug")
	var record domain.Record
	if err := c.Bind().Body(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	id, err := h.service.CreateRecord(ctx, slug, record)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

// GetRecord find a single record by ID
// @Summary Get record by ID
// @Description Retrieve a specific dynamic record by its unique MongoDB ID
// @Tags content
// @Produce json
// @Param collSlug path string true "Collection Slug"
// @Param id path string true "Record ID (Hex)"
// @Success 200 {object} map[string]interface{}
// @Router /content/{collSlug}/{id} [get]
func (h *Handler) GetRecord(c fiber.Ctx) error {
	slug := c.Params("slug")
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid object id"})
	}

	ctx := context.Background()
	record, err := h.service.GetRecord(ctx, slug, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "record not found"})
	}
	return c.JSON(record)
}

// UpdateRecord update an existing record
// @Summary Update record by ID
// @Description Modify a dynamic record. Validates update against collection schema.
// @Tags content
// @Accept json
// @Produce json
// @Param collSlug path string true "Collection Slug"
// @Param id path string true "Record ID (Hex)"
// @Param record body map[string]interface{} true "Updated Record Data"
// @Success 204 "No Content"
// @Router /content/{collSlug}/{id} [put]
func (h *Handler) UpdateRecord(c fiber.Ctx) error {
	slug := c.Params("slug")
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
	if err := h.service.UpdateRecord(ctx, slug, id, record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteRecord delete a record by ID
// @Summary Delete record by ID
// @Description Remove a specific dynamic record from a collection
// @Tags content
// @Param collSlug path string true "Collection Slug"
// @Param id path string true "Record ID (Hex)"
// @Success 204 "No Content"
// @Router /content/{collSlug}/{id} [delete]
func (h *Handler) DeleteRecord(c fiber.Ctx) error {
	slug := c.Params("slug")
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid object id"})
	}

	ctx := context.Background()
	if err := h.service.DeleteRecord(ctx, slug, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
