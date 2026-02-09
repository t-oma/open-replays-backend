package validation

import (
	"fmt"
	"strings"

	"open-replays/internal/api/response"
)

// StringValidator validates string values using fluent API.
type StringValidator struct {
	base *baseValidator[string]
}

// String creates a new string validator.
func String(field, value string) *StringValidator {
	return &StringValidator{
		base: newBaseValidator(field, value),
	}
}

// When conditionally applies validation only if the condition is true.
func (v *StringValidator) When(condition bool) *StringValidator {
	v.base.When(condition)
	return v
}

// Optional skips validation if value is empty (zero value).
func (v *StringValidator) Optional() *StringValidator {
	return v.When(v.base.Value != "")
}

// Required checks if value is not empty.
func (v *StringValidator) Required() *StringValidator {
	v.base.Required(strings.TrimSpace(v.base.Value) != "")
	return v
}

// MinLength checks minimum length.
func (v *StringValidator) MinLength(minLen int) *StringValidator {
	if v.base.SkipValidation {
		return v
	}
	if len(v.base.Value) < minLen {
		v.base.AddError(
			fmt.Sprintf("%s must be at least %d characters", v.base.Field, minLen),
			response.TagMinLength,
			response.LengthDetails{
				MinLength:    minLen,
				ActualLength: len(v.base.Value),
				ExceededBy:   minLen - len(v.base.Value),
			},
		)
	}
	return v
}

// MaxLength checks maximum length.
func (v *StringValidator) MaxLength(maxLen int) *StringValidator {
	if v.base.SkipValidation {
		return v
	}
	if len(v.base.Value) > maxLen {
		v.base.AddError(
			fmt.Sprintf("%s must not exceed %d characters", v.base.Field, maxLen),
			response.TagMaxLength,
			response.LengthDetails{
				MaxLength:    maxLen,
				ActualLength: len(v.base.Value),
				ExceededBy:   len(v.base.Value) - maxLen,
			},
		)
	}
	return v
}

// Collect returns collected errors.
func (v *StringValidator) Collect() []response.ValidationError {
	return v.base.Collect()
}

// IsValid returns true if no errors.
func (v *StringValidator) IsValid() bool {
	return v.base.IsValid()
}
