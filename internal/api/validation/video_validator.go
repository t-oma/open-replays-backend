package validation

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"open-replays/internal/api/response"
)

type VideoValidatorImpl struct {
	rules VideoRules
}

var _ VideoValidator = (*VideoValidatorImpl)(nil)

// VideoRules contains video validation configuration.
type VideoRules struct {
	MinTitleLength int
	MaxTitleLength int

	MaxDescriptionLength int

	MinVideoSize int64
	MaxVideoSize int64

	MinThumbnailSize int64
	MaxThumbnailSize int64

	AllowedVideoExts     []string
	AllowedThumbnailExts []string
}

// NewVideoValidator creates a new VideoValidator with default rules.
func NewVideoValidator() VideoValidator {
	return &VideoValidatorImpl{
		rules: DefaultValidationRules(),
	}
}

// NewVideoValidatorWithRules creates validator with custom rules.
func NewVideoValidatorWithRules(rules VideoRules) VideoValidator {
	return &VideoValidatorImpl{rules: rules}
}

// DefaultValidationRules returns default validation rules.
func DefaultValidationRules() VideoRules {
	return VideoRules{
		MinTitleLength: MinTitleLength,
		MaxTitleLength: MaxTitleLength,

		MaxDescriptionLength: MaxDescriptionLength,

		MinVideoSize: MinVideoSize,
		MaxVideoSize: MaxVideoSize,

		MinThumbnailSize: MinThumbnailSize,
		MaxThumbnailSize: MaxThumbnailSize,

		AllowedVideoExts:     []string{".mp4", ".mov", ".webm"},
		AllowedThumbnailExts: []string{".jpg", ".jpeg", ".png"},
	}
}

func (v *VideoValidatorImpl) ValidateUpload(req UploadVideoRequest) []response.ValidationError {
	var errors []response.ValidationError

	// Title
	errors = append(errors, v.ValidateTitle(req.GetTitle())...)

	description := req.GetDescription()
	if description != "" && len(description) > v.rules.MaxDescriptionLength {
		errors = append(errors, response.ValidationError{
			Field:   "description",
			Message: "Description must not exceed 2000 characters",
			Tag:     TagMaxLength,
			Details: response.LengthDetails{
				MaxLength:    v.rules.MaxDescriptionLength,
				ActualLength: len(description),
				ExceededBy:   len(description) - v.rules.MaxDescriptionLength,
			},
		})
	}

	// Video
	errors = append(errors, v.ValidateVideoFile(req.GetVideo())...)
	// Thumbnail
	errors = append(errors, v.ValidateThumbnailFile(req.GetThumbnail())...)

	return errors
}

func (v *VideoValidatorImpl) ValidateTitle(title string) []response.ValidationError {
	var errors []response.ValidationError

	switch {
	case strings.TrimSpace(title) == "":
		errors = append(errors, response.ValidationError{
			Field:   "title",
			Message: "Title is required",
			Tag:     TagRequired,
		})
	case len(title) > v.rules.MaxTitleLength:
		errors = append(errors, response.ValidationError{
			Field:   "title",
			Message: "Title must not exceed 200 characters",
			Tag:     TagMaxLength,
			Details: response.LengthDetails{
				MaxLength:    v.rules.MaxTitleLength,
				ActualLength: len(title),
				ExceededBy:   len(title) - v.rules.MaxTitleLength,
			},
		})
	case len(title) < v.rules.MinTitleLength:
		errors = append(errors, response.ValidationError{
			Field:   "title",
			Message: "Title must be at least 5 characters",
			Tag:     TagMinLength,
			Details: response.LengthDetails{
				MinLength:    v.rules.MinTitleLength,
				ActualLength: len(title),
				ExceededBy:   len(title) - v.rules.MinTitleLength,
			},
		})
	}

	return errors
}

func (v *VideoValidatorImpl) ValidateVideoFile(
	file *multipart.FileHeader,
) []response.ValidationError {
	var errors []response.ValidationError

	switch {
	case file == nil:
		errors = append(errors, response.ValidationError{
			Field:   "video",
			Message: "Video file is required",
			Tag:     TagRequired,
		})
		return errors
	case file.Size > v.rules.MaxVideoSize:
		errors = append(errors, response.ValidationError{
			Field:   "video",
			Message: "Video file size must not exceed 100MB",
			Tag:     TagMaxSize,
			Details: response.FileSizeDetails{
				MaxSizeBytes:    v.rules.MaxVideoSize,
				MaxSizeMB:       float64(v.rules.MaxVideoSize) / (1024 * 1024),
				ActualSizeBytes: file.Size,
				ActualSizeMB:    float64(file.Size) / (1024 * 1024),
			},
		})
	case file.Size < v.rules.MinVideoSize:
		errors = append(errors, response.ValidationError{
			Field:   "video",
			Message: "Video file size must be at least 1MB",
			Tag:     TagMinSize,
			Details: response.FileSizeDetails{
				MinSizeBytes:    v.rules.MinVideoSize,
				MinSizeMB:       float64(v.rules.MinVideoSize) / (1024 * 1024),
				ActualSizeBytes: file.Size,
				ActualSizeMB:    float64(file.Size) / (1024 * 1024),
			},
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !IsVideoExtensionAllowed(ext) {
		errors = append(errors, response.ValidationError{
			Field:   "video",
			Message: "Video must be in MP4, MOV, or WEBM format",
			Tag:     TagInvalidFormat,
			Details: response.FileTypeDetails{
				AllowedTypes: v.rules.AllowedVideoExts,
				ActualType:   ext,
				Filename:     file.Filename,
			},
		})
	}

	return errors
}

func (v *VideoValidatorImpl) ValidateThumbnailFile(
	file *multipart.FileHeader,
) []response.ValidationError {
	var errors []response.ValidationError

	switch {
	case file == nil || file.Size <= 0:
		return errors
	case file.Size > v.rules.MaxThumbnailSize:
		errors = append(errors, response.ValidationError{
			Field:   "thumbnail",
			Message: "Thumbnail size must not exceed 5MB",
			Tag:     TagMaxSize,
			Details: response.FileSizeDetails{
				MaxSizeBytes:    v.rules.MaxThumbnailSize,
				MaxSizeMB:       float64(v.rules.MaxThumbnailSize) / (1024 * 1024),
				ActualSizeBytes: file.Size,
				ActualSizeMB:    float64(file.Size) / (1024 * 1024),
			},
		})
	case file.Size < v.rules.MinThumbnailSize:
		errors = append(errors, response.ValidationError{
			Field:   "thumbnail",
			Message: "Thumbnail size must be at least 1KB",
			Tag:     TagMinSize,
			Details: response.FileSizeDetails{
				MinSizeBytes:    v.rules.MinThumbnailSize,
				MinSizeMB:       float64(v.rules.MinThumbnailSize) / (1024 * 1024),
				ActualSizeBytes: file.Size,
				ActualSizeMB:    float64(file.Size) / (1024 * 1024),
			},
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !IsThumbnailExtensionAllowed(ext) {
		errors = append(errors, response.ValidationError{
			Field:   "thumbnail",
			Message: "Thumbnail must be in JPG, JPEG, or PNG format",
			Tag:     TagInvalidFormat,
			Details: response.FileTypeDetails{
				AllowedTypes: v.rules.AllowedThumbnailExts,
				ActualType:   ext,
				Filename:     file.Filename,
			},
		})
	}

	return errors
}
