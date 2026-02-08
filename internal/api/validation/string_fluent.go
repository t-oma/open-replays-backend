package validation

import (
	"fmt"
	"strings"

	"open-replays/internal/api/response"
)

// StringValidator validates string values using fluent API.
type StringValidator struct {
	field          string
	value          string
	errors         []response.ValidationError
	skipValidation bool
}

// String creates a new string validator.
func String(field, value string) *StringValidator {
	return &StringValidator{
		field:  field,
		value:  value,
		errors: make([]response.ValidationError, 0),
	}
}

// When conditionally applies validation only if the condition is true.
func (v *StringValidator) When(condition bool) *StringValidator {
	v.skipValidation = !condition
	return v
}

// Optional skips validation if value is empty (zero value).
func (v *StringValidator) Optional() *StringValidator {
	return v.When(v.value != "")
}

// Required checks if value is not empty.
func (v *StringValidator) Required() *StringValidator {
	if v.skipValidation {
		return v
	}
	if strings.TrimSpace(v.value) == "" {
		v.errors = append(v.errors, response.ValidationError{
			Field:   v.field,
			Message: v.field + " is required",
			Tag:     TagRequired,
		})
	}
	return v
}

// MinLength checks minimum length.
func (v *StringValidator) MinLength(minLen int) *StringValidator {
	if v.skipValidation {
		return v
	}
	if len(v.value) < minLen {
		v.errors = append(v.errors, response.ValidationError{
			Field:   v.field,
			Message: fmt.Sprintf("%s must be at least %d characters", v.field, minLen),
			Tag:     TagMinLength,
			Details: response.LengthDetails{
				MinLength:    minLen,
				ActualLength: len(v.value),
				ExceededBy:   minLen - len(v.value),
			},
		})
	}
	return v
}

// MaxLength checks maximum length.
func (v *StringValidator) MaxLength(maxLen int) *StringValidator {
	if v.skipValidation {
		return v
	}
	if len(v.value) > maxLen {
		v.errors = append(v.errors, response.ValidationError{
			Field:   v.field,
			Message: fmt.Sprintf("%s must not exceed %d characters", v.field, maxLen),
			Tag:     TagMaxLength,
			Details: response.LengthDetails{
				MaxLength:    maxLen,
				ActualLength: len(v.value),
				ExceededBy:   len(v.value) - maxLen,
			},
		})
	}
	return v
}

// Errors returns collected errors.
func (v *StringValidator) Errors() []response.ValidationError {
	return v.errors
}

// IsValid returns true if no errors.
func (v *StringValidator) IsValid() bool {
	return len(v.errors) == 0
}
