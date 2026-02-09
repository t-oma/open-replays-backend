package validation

import "open-replays/internal/api/response"

type baseValidator[T any] struct {
	Field          string
	Value          T
	Errors         []response.ValidationError
	SkipValidation bool
}

func newBaseValidator[T any](field string, value T) *baseValidator[T] {
	return &baseValidator[T]{
		Field:          field,
		Value:          value,
		Errors:         make([]response.ValidationError, 0),
		SkipValidation: false,
	}
}

// When conditionally applies validation only if the condition is true.
func (v *baseValidator[T]) When(condition bool) *baseValidator[T] {
	v.SkipValidation = !condition
	return v
}

// Required checks if value is not empty.
func (v *baseValidator[T]) Required(condition bool) *baseValidator[T] {
	if v.SkipValidation {
		return v
	}
	if !condition {
		v.AddError(v.Field+" is required", response.TagRequired, nil)
	}
	return v
}

func (v *baseValidator[T]) AddError(
	message string,
	tag response.ValidationTag,
	details any,
) *baseValidator[T] {
	v.Errors = append(v.Errors, response.ValidationError{
		Field:   v.Field,
		Message: message,
		Tag:     tag,
		Details: details,
	})
	return v
}

// Errors returns collected errors.
func (v *baseValidator[T]) Collect() []response.ValidationError {
	return v.Errors
}

// IsValid returns true if no errors.
func (v *baseValidator[T]) IsValid() bool {
	return len(v.Errors) == 0
}

// FirstError returns first validation error if any.
func (v *baseValidator[T]) FirstError() *response.ValidationError {
	if len(v.Errors) > 0 {
		return &v.Errors[0]
	}
	return nil
}
