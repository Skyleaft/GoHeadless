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
	group.Post("/", h.UploadFile)
	group.Delete("/", h.DeleteFile)
	group.Put("/", h.UpdateFile)
}

// UploadFile handles generic file uploads
// @Summary Upload a file
// @Description Handles file upload, converts images to WebP if applicable, and saves it to the server's local storage.
// @Tags Upload
// @Accept multipart/form-data
// @Param file formData file true "The file to upload"
// @Success 200 {object} map[string]string "Returns the file path"
// @Failure 400 {object} map[string]string "Invalid file upload request"
// @Failure 500 {object} map[string]string "Server error during processing"
// @Router /upload [post]
func (h *Handler) UploadFile(c fiber.Ctx) error {
	// 1. Parse Multipart Form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to get file from form: %v", err),
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
	filePath, err := h.service.UploadFile(file, fileHeader.Filename)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to process file: %v", err),
		})
	}

	// 4. Return the file path
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"path": filePath,
	})
}

// DeleteFile handles deleting an uploaded file
// @Summary Delete an uploaded file
// @Description Deletes a file from the server's local storage based on its path.
// @Tags Upload
// @Param path query string true "The path of the file to delete (e.g. /uploads/file_name.ext)"
// @Success 200 {object} map[string]string "Returns success message"
// @Failure 400 {object} map[string]string "Invalid path or file not found"
// @Failure 500 {object} map[string]string "Server error during file deletion"
// @Router /upload [delete]
func (h *Handler) DeleteFile(c fiber.Ctx) error {
	filePath := c.Query("path")
	if filePath == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "path query parameter is required",
		})
	}

	if err := h.service.DeleteFile(filePath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to delete file: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "file deleted successfully",
	})
}

// UpdateFile handles replacing an existing file
// @Summary Replace an uploaded file
// @Description Deletes the old file and replaces it with a new one.
// @Tags Upload
// @Accept multipart/form-data
// @Param oldPath query string true "The path of the file to replace"
// @Param file formData file true "The new file"
// @Success 200 {object} map[string]string "Returns the new file path"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error during processing"
// @Router /upload [put]
func (h *Handler) UpdateFile(c fiber.Ctx) error {
	oldPath := c.Query("oldPath")
	if oldPath == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "oldPath query parameter is required",
		})
	}

	// 1. Parse Multipart Form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to get file from form: %v", err),
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
	filePath, err := h.service.UpdateFile(oldPath, file, fileHeader.Filename)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to process file update: %v", err),
		})
	}

	// 4. Return the new path
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"path": filePath,
	})
}
