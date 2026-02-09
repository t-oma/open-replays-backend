package validation

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
