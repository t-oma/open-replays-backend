// Package validator provides HTTP request validation.
package validator

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"open-replays/internal/api/httperr"
)

// ValidationError represents a single validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult holds validation errors.
type ValidationResult struct {
	Errors []ValidationError `json:"errors"`
}

// Add adds a validation error.
func (v *ValidationResult) Add(field, message string) {
	v.Errors = append(v.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// IsValid returns true if there are no validation errors.
func (v *ValidationResult) IsValid() bool {
	return len(v.Errors) == 0
}

// ToError converts validation result to AppError.
func (v *ValidationResult) ToError() *httperr.AppError {
	return httperr.ErrValidation.WithDetails(v.Errors)
}

// UploadValidator validates upload requests.
type UploadValidator struct {
	MaxFileSize       int64
	AllowedExtensions map[string]bool
}

// NewUploadValidator creates a new upload validator.
func NewUploadValidator(maxFileSize int64, allowedExtensions []string) *UploadValidator {
	extMap := make(map[string]bool, len(allowedExtensions))
	for _, ext := range allowedExtensions {
		extMap[ext] = true
	}
	return &UploadValidator{
		MaxFileSize:       maxFileSize,
		AllowedExtensions: extMap,
	}
}

// ValidateFile validates a file upload.
func (v *UploadValidator) ValidateFile(file *multipart.FileHeader) *httperr.AppError {
	if file == nil {
		return httperr.ErrInvalidRequest.WithDetails(map[string]string{
			"field": "file",
			"value": "nil",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !v.AllowedExtensions[ext] {
		return httperr.ErrInvalidFileType.WithDetails(FileTypeDetails{
			AllowedTypes: v.getAllowedExtensionsList(),
			ActualType:   ext,
			Filename:     file.Filename,
		})
	}

	if file.Size > v.MaxFileSize {
		return httperr.ErrFileTooLarge.WithDetails(FileSizeDetails{
			MaxSizeBytes:    v.MaxFileSize,
			MaxSizeMB:       float64(v.MaxFileSize) / (1024 * 1024),
			ActualSizeBytes: file.Size,
			ActualSizeMB:    float64(file.Size) / (1024 * 1024),
		})
	}

	return nil
}

// ValidateTitle validates video title.
func (v *UploadValidator) ValidateTitle(title string) *ValidationResult {
	result := &ValidationResult{}

	if strings.TrimSpace(title) == "" {
		result.Add("title", "title is required")
	}

	if len(title) > 200 {
		result.Add("title", "title must not exceed 200 characters")
	}

	return result
}

// FileTypeDetails provides details for file type errors.
type FileTypeDetails struct {
	AllowedTypes []string `json:"allowedTypes"`
	ActualType   string   `json:"actualType"`
	Filename     string   `json:"filename"`
}

// FileSizeDetails provides details for file size errors.
type FileSizeDetails struct {
	MaxSizeBytes    int64   `json:"maxSizeBytes"`
	MaxSizeMB       float64 `json:"maxSizeMb"`
	ActualSizeBytes int64   `json:"actualSizeBytes"`
	ActualSizeMB    float64 `json:"actualSizeMb"`
}

// getAllowedExtensionsList returns a list of allowed extensions.
func (v *UploadValidator) getAllowedExtensionsList() []string {
	exts := make([]string, 0, len(v.AllowedExtensions))
	for ext := range v.AllowedExtensions {
		exts = append(exts, ext)
	}
	return exts
}
