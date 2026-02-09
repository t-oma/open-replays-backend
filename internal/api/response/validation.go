package response

// ValidationTag represents the type of validation error.
type ValidationTag string

// Validation error tags.
const (
	TagRequired          ValidationTag = "required"
	TagMinLength         ValidationTag = "min_length"
	TagMaxLength         ValidationTag = "max_length"
	TagMinFileSize       ValidationTag = "min_file_size"
	TagMaxFileSize       ValidationTag = "max_file_size"
	TagMinIntSize        ValidationTag = "min_int_size"
	TagMaxIntSize        ValidationTag = "max_int_size"
	TagInvalidFileFormat ValidationTag = "invalid_file_format"
)

// ValidationError represents validation error with optional details.
type ValidationError struct {
	Field   string        `json:"field"`
	Message string        `json:"message"`
	Tag     ValidationTag `json:"tag"`
	Details any           `json:"details,omitempty"`
}

// FileTypeDetails contains file type validation details.
type FileTypeDetails struct {
	AllowedTypes []string `json:"allowedTypes"`
	ActualType   string   `json:"actualType"`
	Filename     string   `json:"filename"`
}

// FileSizeDetails contains file size validation details.
type FileSizeDetails struct {
	MinSizeBytes    int64   `json:"minSizeBytes,omitempty"`
	MaxSizeBytes    int64   `json:"maxSizeBytes,omitempty"`
	MinSizeMB       float64 `json:"minSizeMb,omitempty"`
	MaxSizeMB       float64 `json:"maxSizeMb,omitempty"`
	ActualSizeBytes int64   `json:"actualSizeBytes"`
	ActualSizeMB    float64 `json:"actualSizeMb"`
}

// LengthDetails contains string length validation details.
type LengthDetails struct {
	MinLength    int `json:"minLength,omitempty"`
	MaxLength    int `json:"maxLength,omitempty"`
	ActualLength int `json:"actualLength"`
	ExceededBy   int `json:"exceededBy"`
}

// IntSizeDetails contains int64 size validation details.
type IntSizeDetails struct {
	MinSize    int64 `json:"minSize,omitempty"`
	MaxSize    int64 `json:"maxSize,omitempty"`
	ActualSize int64 `json:"actualSize"`
	ExceededBy int64 `json:"exceededBy"`
}
