package upload

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nickalie/go-webpbin"
)

type Service interface {
	UploadFile(file io.Reader, filename string) (string, error)
	DeleteFile(path string) error
	UpdateFile(oldPath string, file io.Reader, filename string) (string, error)
}

type service struct {
	basePath string
}

func NewService(basePath string) Service {
	return &service{basePath: basePath}
}

func (s *service) UploadFile(file io.Reader, filename string) (string, error) {
	// 1. Read entire file into memory (since these are typical uploads)
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read upload data: %w", err)
	}

	contentType := http.DetectContentType(data)
	ext := filepath.Ext(filename)
	cleanName := strings.TrimSuffix(filename, ext)
	timestamp := time.Now().UnixNano()

	// 2. If it's an image, attempt WebP conversion
	if strings.HasPrefix(contentType, "image/") && !strings.Contains(contentType, "svg") {
		img, _, err := image.Decode(bytes.NewReader(data))
		if err == nil {
			newName := fmt.Sprintf("%s_%d.webp", cleanName, timestamp)
			if err := s.ensureDir(); err != nil {
				return "", err
			}
			filePath := filepath.Join(s.basePath, newName)

			if err := webpbin.NewCWebP().
				Quality(80).
				InputImage(img).
				OutputFile(filePath).
				Run(); err == nil {
				return filepath.Join("/uploads", newName), nil
			}
			// Fallback if conversion fails
		}
	}

	// 3. Save as-is (fallback or non-image)
	newName := fmt.Sprintf("%s_%d%s", cleanName, timestamp, ext)
	filePath := filepath.Join(s.basePath, newName)
	if err := s.ensureDir(); err != nil {
		return "", err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filepath.Join("/uploads", newName), nil
}

func (s *service) ensureDir() error {
	if _, err := os.Stat(s.basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(s.basePath, 0755); err != nil {
			return fmt.Errorf("failed to create upload directory: %w", err)
		}
	}
	return nil
}

func (s *service) DeleteFile(filePath string) error {
	if !strings.HasPrefix(filePath, "/uploads") {
		return fmt.Errorf("invalid path: must start with /uploads")
	}

	fileName := strings.TrimPrefix(filePath, "/uploads")
	fileName = strings.TrimPrefix(fileName, "/")
	fullPath := filepath.Join(s.basePath, fileName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *service) UpdateFile(oldPath string, file io.Reader, filename string) (string, error) {
	if err := s.DeleteFile(oldPath); err != nil {
		fmt.Printf("Warning: failed to delete old file %s during update: %v\n", oldPath, err)
	}
	return s.UploadFile(file, filename)
}
