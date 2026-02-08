package response

// ValidationError represents validation error with optional details.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Details any    `json:"details,omitempty"`
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
