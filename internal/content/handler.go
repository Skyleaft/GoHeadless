package content

import (
	"context"
	"errors"

	"GoHeadless/internal/apierr"
	"GoHeadless/internal/domain"
	"GoHeadless/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes(router fiber.Router) {
	router.Get("", h.ListRecords)
	router.Post("", h.CreateRecord)
	router.Get("/:id", h.GetRecord)
	router.Put("/:id", h.UpdateRecord)
	router.Delete("/:id", h.DeleteRecord)
}

// ListRecords lists records with optional search, filter, sort, and pagination.
// @Summary List records in a collection (query engine)
// @Description Paginated list with search, filter[field][op], and sort (-field for DESC). Anonymous users on public collections receive internal fields stripped.
// @Tags content
// @Produce json
// @Param collSlug path string true "Collection Slug"
// @Param page query int false "Page (default 1)"
// @Param limit query int false "Page size (default 10, max 100)"
// @Param search query string false "Search across searchable text fields"
// @Param sort query string false "Sort field; prefix - for descending"
// @Param filter query string false "Dynamic filters: filter[key]=val or filter[key][gt]=val"
// @Success 200 {object} ListRecordsResult
// @Router /content/{collSlug} [get]
func (h *Handler) ListRecords(c fiber.Ctx) error {
	slug := c.Params("slug")
	parser := NewQueryParser()
	pq, err := parser.Parse(string(c.Request().URI().QueryString()))
	if err != nil {
		return apierr.BadRequest(err.Error())
	}

	if c.Locals(middleware.CollectionUnknownLocalsKey) == true {
		return c.JSON(ListRecordsResult{
			Data:  []domain.Record{},
			Total: 0,
			Page:  pq.Page,
			Limit: pq.Limit,
		})
	}

	stripInternal := c.Locals("user_id") == nil

	out, err := h.service.ListRecords(context.Background(), slug, pq, stripInternal)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidQuery):
			return apierr.BadRequest(err.Error())
		case errors.Is(err, ErrCollectionNotFound):
			return apierr.NotFound("collection not found")
		default:
			return apierr.Internal("failed to list records", err)
		}
	}
	return c.JSON(out)
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
		return apierr.BadRequest("invalid request body")
	}

	id, err := h.service.CreateRecord(context.Background(), slug, record)
	if err != nil {
		return apierr.BadRequest(err.Error())
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
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return apierr.BadRequest("invalid object id")
	}

	record, err := h.service.GetRecord(context.Background(), slug, id)
	if err != nil {
		return apierr.NotFound("record not found")
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
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return apierr.BadRequest("invalid object id")
	}

	var record domain.Record
	if err := c.Bind().Body(&record); err != nil {
		return apierr.BadRequest("invalid request body")
	}

	if err := h.service.UpdateRecord(context.Background(), slug, id, record); err != nil {
		return apierr.BadRequest(err.Error())
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
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return apierr.BadRequest("invalid object id")
	}

	if err := h.service.DeleteRecord(context.Background(), slug, id); err != nil {
		return apierr.Internal("failed to delete record", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
