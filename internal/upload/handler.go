package upload

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes(router fiber.Router) {
	group := router.Group("/upload")
	group.Post("/image", h.UploadImage)
	group.Delete("/image", h.DeleteImage)
	group.Put("/image", h.UpdateImage)
}

// UploadImage handles image uploads and converts them to webp
// @Summary Upload and convert image to webp
// @Description Handles image file upload, converts it to WebP format, and saves it to the server's local storage.
// @Tags Upload
// @Accept multipart/form-data
// @Param image formData file true "The image to upload"
// @Success 200 {object} map[string]string "Returns the image path"
// @Failure 400 {object} map[string]string "Invalid file upload request"
// @Failure 500 {object} map[string]string "Server error during image processing"
// @Router /upload/image [post]
func (h *Handler) UploadImage(c fiber.Ctx) error {
	// 1. Parse Multipart Form
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to get image from form: %v", err),
		})
	}

	// 2. Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to open uploaded file: %v", err),
		})
	}
	defer file.Close()

	// 3. Process the file through the service
	imagePath, err := h.service.UploadImage(file, fileHeader.Filename)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to process image: %v", err),
		})
	}

	// 4. Return the image path for subsequent storage in CMS
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"path": imagePath,
	})
}

// DeleteImage handles deleting an uploaded image
// @Summary Delete an uploaded image
// @Description Deletes an image file from the server's local storage based on its path.
// @Tags Upload
// @Param path query string true "The path of the image to delete (e.g. /uploads/image_name.webp)"
// @Success 200 {object} map[string]string "Returns success message"
// @Failure 400 {object} map[string]string "Invalid path or file not found"
// @Failure 500 {object} map[string]string "Server error during file deletion"
// @Router /upload/image [delete]
func (h *Handler) DeleteImage(c fiber.Ctx) error {
	imagePath := c.Query("path")
	if imagePath == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "path query parameter is required",
		})
	}

	if err := h.service.DeleteImage(imagePath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to delete image: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "image deleted successfully",
	})
}

// UpdateImage handles replacing an existing image
// @Summary Replace an uploaded image
// @Description Deletes the old image and replaces it with a new one.
// @Tags Upload
// @Accept multipart/form-data
// @Param oldPath query string true "The path of the image to replace"
// @Param image formData file true "The new image file"
// @Success 200 {object} map[string]string "Returns the new image path"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error during processing"
// @Router /upload/image [put]
func (h *Handler) UpdateImage(c fiber.Ctx) error {
	oldPath := c.Query("oldPath")
	if oldPath == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "oldPath query parameter is required",
		})
	}

	// 1. Parse Multipart Form
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to get image from form: %v", err),
		})
	}

	// 2. Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to open uploaded file: %v", err),
		})
	}
	defer file.Close()

	// 3. Process the file through the service
	imagePath, err := h.service.UpdateImage(oldPath, file, fileHeader.Filename)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to process image update: %v", err),
		})
	}

	// 4. Return the new image path
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"path": imagePath,
	})
}
