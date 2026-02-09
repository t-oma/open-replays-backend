package validation

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"

	"open-replays/internal/api/response"
)

// FileValidator validates file uploads using fluent API.
type FileValidator struct {
	base *baseValidator[*multipart.FileHeader]
}

// File creates a new file validator.
func File(field string, file *multipart.FileHeader) *FileValidator {
	return &FileValidator{
		base: newBaseValidator(field, file),
	}
}

// When conditionally applies validation only if the condition is true.
func (v *FileValidator) When(condition bool) *FileValidator {
	v.base.When(condition)
	return v
}

// Optional skips validation if file is nil or empty.
func (v *FileValidator) Optional() *FileValidator {
	return v.When(v.base.Value != nil && v.base.Value.Size > 0)
}

// Required checks if file is provided.
func (v *FileValidator) Required() *FileValidator {
	v.base.Required(v.base.Value != nil && v.base.Value.Size > 0)
	return v
}

// MinSize checks minimum file size.
func (v *FileValidator) MinSize(minLen int64) *FileValidator {
	if v.base.SkipValidation {
		return v
	}
	if v.base.Value != nil && v.base.Value.Size < minLen {
		v.base.AddError(
			v.base.Field+" is too small",
			response.TagMinFileSize,
			response.FileSizeDetails{
				MinSizeBytes:    minLen,
				ActualSizeBytes: v.base.Value.Size,
			},
		)
	}
	return v
}

// MaxSize checks maximum file size.
func (v *FileValidator) MaxSize(maxLen int64) *FileValidator {
	if v.base.SkipValidation {
		return v
	}
	if v.base.Value != nil && v.base.Value.Size > maxLen {
		v.base.AddError(
			fmt.Sprintf("%s must not exceed %s", v.base.Field, formatSize(maxLen)),
			response.TagMaxFileSize,
			response.FileSizeDetails{
				MaxSizeBytes:    maxLen,
				ActualSizeBytes: v.base.Value.Size,
			},
		)
	}
	return v
}

// AllowedExts checks if file extension is allowed.
func (v *FileValidator) AllowedExts(exts []string) *FileValidator {
	if v.base.SkipValidation {
		return v
	}
	if v.base.Value == nil {
		return v
	}

	ext := strings.ToLower(filepath.Ext(v.base.Value.Filename))
	if slices.Contains(exts, ext) {
		return v
	}

	v.base.AddError(
		v.base.Field+" has invalid format",
		response.TagInvalidFileFormat,
		response.FileTypeDetails{
			AllowedTypes: exts,
			ActualType:   ext,
			Filename:     v.base.Value.Filename,
		},
	)
	return v
}

// Collect returns collected errors.
func (v *FileValidator) Collect() []response.ValidationError {
	return v.base.Collect()
}

// IsValid returns true if no errors.
func (v *FileValidator) IsValid() bool {
	return v.base.IsValid()
}

// formatSize formats bytes to human readable string.
func formatSize(bytes int64) string {
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}
