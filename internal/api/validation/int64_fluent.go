package validation

import (
	"open-replays/internal/api/response"
)

// Int64Validator validates int64 values using fluent API.
type Int64Validator struct {
	base *baseValidator[int64]
}

// Int64 creates a new int64 validator.
func Int64(field string, value int64) *Int64Validator {
	return &Int64Validator{
		base: newBaseValidator(field, value),
	}
}

// Min checks minimum value.
func (v *Int64Validator) Min(minSize int64) *Int64Validator {
	if v.base.Value < minSize {
		v.base.AddError(
			v.base.Field+" is too small",
			response.TagMinIntSize,
			response.IntSizeDetails{
				MinSize:    minSize,
				ActualSize: v.base.Value,
			},
		)
	}
	return v
}

// Max checks maximum value.
func (v *Int64Validator) Max(maxSize int64) *Int64Validator {
	if v.base.Value > maxSize {
		v.base.AddError(
			v.base.Field+" is too large",
			response.TagMaxIntSize,
			response.IntSizeDetails{
				MaxSize:    maxSize,
				ActualSize: v.base.Value,
			},
		)
	}
	return v
}

// Range checks if value is within range [minSize, maxSize].
func (v *Int64Validator) Range(minSize, maxSize int64) *Int64Validator {
	return v.Min(minSize).Max(maxSize)
}

// Collect returns collected errors.
func (v *Int64Validator) Collect() []response.ValidationError {
	return v.base.Collect()
}

// IsValid returns true if no errors.
func (v *Int64Validator) IsValid() bool {
	return v.base.IsValid()
}
