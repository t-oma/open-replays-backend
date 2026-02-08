package validation

import "slices"

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

const (
	MinTitleLength       = 5
	MaxTitleLength       = 100
	MaxDescriptionLength = 2000

	MinVideoSize     = KB
	MaxVideoSize     = 100 * MB
	MinThumbnailSize = KB
	MaxThumbnailSize = 5 * MB
)

var (
	AllowedVideoExts     = []string{".mp4", ".mov", ".webm"}
	AllowedThumbnailExts = []string{".jpg", ".jpeg", ".png"}
)

func IsVideoExtensionAllowed(ext string) bool {
	return slices.Contains(AllowedVideoExts, ext)
}

func IsThumbnailExtensionAllowed(ext string) bool {
	return slices.Contains(AllowedThumbnailExts, ext)
}
