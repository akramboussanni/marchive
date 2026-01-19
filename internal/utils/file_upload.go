package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// Allowed file extensions
var (
	AllowedBookFormats = []string{".pdf", ".epub", ".mobi", ".azw3", ".djvu", ".fb2", ".txt"}
	AllowedImageFormats = []string{".jpg", ".jpeg", ".png", ".webp"}
)

// File size limits in bytes
const (
	MaxBookSize  = 500 * 1024 * 1024 // 500 MB
	MaxImageSize = 10 * 1024 * 1024  // 10 MB
)

// SaveUploadedFile saves a multipart file to disk with a unique filename
func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, destDir string) (string, error) {
	// Ensure upload directory exists
	if err := EnsureUploadDir(destDir); err != nil {
		return "", err
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	uniqueName := GenerateUniqueFilename(ext)
	destPath := filepath.Join(destDir, uniqueName)

	// Create destination file
	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(destPath) // Clean up on error
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return destPath, nil
}

// ValidateBookFile checks if the book file is valid
func ValidateBookFile(header *multipart.FileHeader) error {
	// Check file size
	if header.Size > MaxBookSize {
		return fmt.Errorf("book file too large: maximum size is %d MB", MaxBookSize/(1024*1024))
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	for _, allowed := range AllowedBookFormats {
		if ext == allowed {
			return nil
		}
	}

	return fmt.Errorf("invalid book format: allowed formats are %v", AllowedBookFormats)
}

// ValidateImageFile checks if the image file is valid
func ValidateImageFile(header *multipart.FileHeader) error {
	// Check file size
	if header.Size > MaxImageSize {
		return fmt.Errorf("image file too large: maximum size is %d MB", MaxImageSize/(1024*1024))
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	for _, allowed := range AllowedImageFormats {
		if ext == allowed {
			return nil
		}
	}

	return fmt.Errorf("invalid image format: allowed formats are %v", AllowedImageFormats)
}

// DeleteFile removes a file from the filesystem
func DeleteFile(path string) error {
	if path == "" {
		return nil
	}
	
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

// GenerateUniqueFilename creates a unique filename using snowflake ID and random string
func GenerateUniqueFilename(ext string) string {
	id := GenerateSnowflakeID()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("%d_%s%s", id, randomStr, ext)
}

// EnsureUploadDir creates the upload directory if it doesn't exist
func EnsureUploadDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}
	return nil
}

// SaveBase64Image saves a base64-encoded image to disk
func SaveBase64Image(base64Data, destDir, ext string) (string, error) {
	// Ensure upload directory exists
	if err := EnsureUploadDir(destDir); err != nil {
		return "", err
	}

	// Remove data URI prefix if present (e.g., "data:image/png;base64,")
	if idx := strings.Index(base64Data, ","); idx != -1 {
		base64Data = base64Data[idx+1:]
	}

	// Decode base64
	data, err := hex.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %w", err)
	}

	// Generate unique filename
	uniqueName := GenerateUniqueFilename(ext)
	destPath := filepath.Join(destDir, uniqueName)

	// Write to file
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return destPath, nil
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// FormatFileSize formats a file size in bytes to a human-readable string
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GenerateHash generates a hash for a file based on its path and current timestamp
func GenerateHash(filePath string) string {
	id := GenerateSnowflakeID()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("upload_%d_%s", id, randomStr)
}

// GetFileExtension returns the file extension from a filename
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 0 && ext[0] == '.' {
		return ext[1:] // Remove the leading dot
	}
	return ext
}

