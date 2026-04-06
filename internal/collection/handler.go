package collection

import (
	"context"

	"GoHeadless/internal/apierr"
	"GoHeadless/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes(router fiber.Router) {
	router.Get("/", h.ListCollections)
	router.Post("/", h.CreateCollection)
	router.Get("/:slug", h.GetCollection)
	router.Put("/:slug", h.UpdateCollection)
	router.Delete("/:slug", h.DeleteCollection)
}

// ListCollections list all collections
// @Summary List all collections
// @Description Get metadata for all user-defined collections
// @Tags collections
// @Produce json
// @Success 200 {array} domain.Collection
// @Router /collections [get]
func (h *Handler) ListCollections(c fiber.Ctx) error {
	colls, err := h.service.GetAllCollections(context.Background())
	if err != nil {
		return apierr.Internal("failed to list collections", err)
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
		return apierr.BadRequest("invalid request body")
	}

	if err := h.service.CreateCollection(context.Background(), &coll); err != nil {
		return apierr.BadRequest(err.Error())
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
	coll, err := h.service.GetCollection(context.Background(), c.Params("slug"))
	if err != nil {
		return apierr.NotFound("collection not found")
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
	if err := h.service.DeleteCollection(context.Background(), c.Params("slug")); err != nil {
		return apierr.Internal("failed to delete collection", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateCollection update a collection by slug
// @Summary Update collection by slug
// @Description Update metadata, fields, and access control for a specific collection
// @Tags collections
// @Accept json
// @Produce json
// @Param slug path string true "Collection Slug"
// @Param collection body domain.Collection true "Collection Definition"
// @Success 200 {object} domain.Collection
// @Router /collections/{slug} [put]
func (h *Handler) UpdateCollection(c fiber.Ctx) error {
	var coll domain.Collection
	if err := c.Bind().Body(&coll); err != nil {
		return apierr.BadRequest("invalid request body")
	}

	if err := h.service.UpdateCollection(context.Background(), c.Params("slug"), &coll); err != nil {
		return apierr.BadRequest(err.Error())
	}
	
	// Return the updated collection object
	updated, err := h.service.GetCollection(context.Background(), c.Params("slug"))
	if err != nil {
		return apierr.Internal("updated successfully but failed to fetch latest response", err)
	}

	return c.JSON(updated)
}
