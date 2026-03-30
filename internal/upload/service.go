package upload

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nickalie/go-webpbin"
)

type Service interface {
	UploadImage(file io.Reader, filename string) (string, error)
	DeleteImage(path string) error
	UpdateImage(oldPath string, file io.Reader, filename string) (string, error)
}

type service struct {
	basePath string
}

func NewService(basePath string) Service {
	return &service{basePath: basePath}
}

func (s *service) UploadImage(file io.Reader, filename string) (string, error) {
	// 1. Decode original image to ensure it's valid
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// 2. Prepare filename (strip extension and add webp)
	ext := filepath.Ext(filename)
	cleanName := strings.TrimSuffix(filename, ext)
	newName := fmt.Sprintf("%s_%d.webp", cleanName, time.Now().UnixNano())
	
	// 3. Ensure upload directory exists
	if _, err := os.Stat(s.basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(s.basePath, 0755); err != nil {
			return "", fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	// 4. Create target file path
	filePath := filepath.Join(s.basePath, newName)

	// 5. Encode into WebP using webpbin (uses pre-compiled binary)
	// We use the NewCWebP() encoder for better cross-platform support without CGO.
	if err := webpbin.NewCWebP().
		Quality(80).
		InputImage(img).
		OutputFile(filePath).
		Run(); err != nil {
		return "", fmt.Errorf("failed to encode into webp using webpbin: %w", err)
	}

	// Returns path starting with /uploads
	return filepath.Join("/uploads", newName), nil
}

func (s *service) DeleteImage(imagePath string) error {
	// 1. Path normalization and security check
	// We only allow deleting files within our /uploads mapped directory.
	if !strings.HasPrefix(imagePath, "/uploads") {
		return fmt.Errorf("invalid path: must start with /uploads")
	}

	// 2. Get the actual file name from the path
	fileName := strings.TrimPrefix(imagePath, "/uploads")
	fileName = strings.TrimPrefix(fileName, "/")

	// 3. Resolve the full physical file path
	fullPath := filepath.Join(s.basePath, fileName)

	// 4. Check if the file exists before attempting deletion
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", imagePath)
	}

	// 5. Delete the file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}

func (s *service) UpdateImage(oldPath string, file io.Reader, filename string) (string, error) {
	// 1. Delete the old image first
	if err := s.DeleteImage(oldPath); err != nil {
		// Log error but proceed if the old image doesn't exist (maybe already deleted)
		fmt.Printf("Warning: failed to delete old image %s during update: %v\n", oldPath, err)
	}

	// 2. Upload the new one
	return s.UploadImage(file, filename)
}
