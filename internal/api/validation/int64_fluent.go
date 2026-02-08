package validation

import (
	"fmt"

	"open-replays/internal/api/response"
)

// Int64Validator validates int64 values using fluent API.
type Int64Validator struct {
	field  string
	value  int64
	errors []response.ValidationError
}

// Int64 creates a new int64 validator.
func Int64(field string, value int64) *Int64Validator {
	return &Int64Validator{
		field:  field,
		value:  value,
		errors: make([]response.ValidationError, 0),
	}
}

// Min checks minimum value.
func (v *Int64Validator) Min(minLen int64) *Int64Validator {
	if v.value < minLen {
		v.errors = append(v.errors, response.ValidationError{
			Field:   v.field,
			Message: fmt.Sprintf("%s is too small", v.field),
			Tag:     TagMinSize,
			Details: response.FileSizeDetails{
				MinSizeBytes:    minLen,
				ActualSizeBytes: v.value,
			},
		})
	}
	return v
}

// Max checks maximum value.
func (v *Int64Validator) Max(maxLen int64) *Int64Validator {
	if v.value > maxLen {
		v.errors = append(v.errors, response.ValidationError{
			Field:   v.field,
			Message: fmt.Sprintf("%s is too large", v.field),
			Tag:     TagMaxSize,
			Details: response.FileSizeDetails{
				MaxSizeBytes:    maxLen,
				ActualSizeBytes: v.value,
			},
		})
	}
	return v
}

// Range checks if value is within range [min, max].
func (v *Int64Validator) Range(min, max int64) *Int64Validator {
	return v.Min(min).Max(max)
}

// Errors returns collected errors.
func (v *Int64Validator) Errors() []response.ValidationError {
	return v.errors
}

// IsValid returns true if no errors.
func (v *Int64Validator) IsValid() bool {
	return len(v.errors) == 0
}
