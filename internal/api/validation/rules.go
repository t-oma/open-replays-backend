// Package validation provides fluent API for request validation.
package validation

import "slices"

// Byte size constants.
const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

// Validation rule constants for videos and thumbnails.
const (
	MinTitleLength       = 5
	MaxTitleLength       = 100
	MaxDescriptionLength = 2000

	MinVideoSize     = KB
	MaxVideoSize     = 100 * MB
	MinThumbnailSize = KB
	MaxThumbnailSize = 5 * MB
)

// Allowed file extensions.
var (
	AllowedVideoExts     = []string{".mp4", ".mov", ".webm"}
	AllowedThumbnailExts = []string{".jpg", ".jpeg", ".png"}
)

// IsVideoExtensionAllowed checks if the given extension is in the allowed video extensions list.
func IsVideoExtensionAllowed(ext string) bool {
	return slices.Contains(AllowedVideoExts, ext)
}

// IsThumbnailExtensionAllowed checks if the given extension is in the allowed thumbnail extensions list.
func IsThumbnailExtensionAllowed(ext string) bool {
	return slices.Contains(AllowedThumbnailExts, ext)
}
