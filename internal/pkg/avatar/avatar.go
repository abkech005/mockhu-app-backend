package avatar

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	MaxFileSize      = 5 * 1024 * 1024 // 5MB
	AvatarSize       = 400              // 400x400 pixels
	StorageDir       = "storage/avatars"
	BaseURL          = "/avatars" // Will be updated for S3 later
)

var (
	ErrFileTooBig       = errors.New("file size exceeds 5MB")
	ErrInvalidFileType  = errors.New("invalid file type, only JPEG, PNG, and WebP allowed")
	ErrProcessingFailed = errors.New("failed to process image")
)

// ProcessAndSave validates, resizes, and saves an avatar image
func ProcessAndSave(fileData []byte, filename string) (string, error) {
	// Validate file size
	if len(fileData) > MaxFileSize {
		return "", ErrFileTooBig
	}

	// Detect and validate file type
	fileType := detectImageType(fileData)
	if fileType == "" {
		return "", ErrInvalidFileType
	}

	// Decode image
	img, err := decodeImage(bytes.NewReader(fileData), fileType)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrProcessingFailed, err)
	}

	// Resize to 400x400 (crop to square, center)
	resized := imaging.Fill(img, AvatarSize, AvatarSize, imaging.Center, imaging.Lanczos)

	// Generate unique filename
	ext := getExtension(fileType)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	
	// Save to local storage (TODO: Add S3 support later)
	_, err = saveToLocal(resized, uniqueFilename, fileType)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return URL (local path for now, will be S3 URL later)
	return fmt.Sprintf("%s/%s", BaseURL, uniqueFilename), nil
}

// DeleteAvatar deletes an avatar file
func DeleteAvatar(avatarURL string) error {
	if avatarURL == "" {
		return nil
	}

	// Extract filename from URL
	filename := filepath.Base(avatarURL)
	if filename == "" || filename == "." {
		return nil
	}

	// Delete from local storage (TODO: Add S3 delete later)
	filePath := filepath.Join(StorageDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	// Delete file
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// detectImageType detects the image type from file data
func detectImageType(data []byte) string {
	// Check PNG
	if len(data) >= 8 && string(data[0:8]) == "\x89PNG\r\n\x1a\n" {
		return "png"
	}
	// Check JPEG
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "jpeg"
	}
	// Check WebP
	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "webp"
	}
	return ""
}

// decodeImage decodes an image based on its type
func decodeImage(r io.Reader, fileType string) (image.Image, error) {
	switch fileType {
	case "png":
		return png.Decode(r)
	case "jpeg":
		return jpeg.Decode(r)
	case "webp":
		// For WebP, use imaging library's generic decode
		return imaging.Decode(r)
	default:
		return nil, ErrInvalidFileType
	}
}

// getExtension returns the file extension for a given type
func getExtension(fileType string) string {
	switch fileType {
	case "png":
		return ".png"
	case "jpeg":
		return ".jpg"
	case "webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

// saveToLocal saves the image to local storage
func saveToLocal(img image.Image, filename, fileType string) (string, error) {
	// Ensure storage directory exists
	if err := os.MkdirAll(StorageDir, 0755); err != nil {
		return "", err
	}

	// Create file path
	filePath := filepath.Join(StorageDir, filename)

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Encode and save (always save as JPEG for consistency and size)
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// TODO: Future S3 implementation
// func saveToS3(img image.Image, filename string) (string, error) {
//     // 1. Encode image to bytes
//     // 2. Upload to S3 bucket
//     // 3. Return S3 URL
//     return "", nil
// }
//
// func deleteFromS3(filename string) error {
//     // 1. Delete from S3 bucket
//     return nil
// }

