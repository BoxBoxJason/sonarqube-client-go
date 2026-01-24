package sonargo

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidValue is returned when a parameter has an invalid value.
	ErrInvalidValue = errors.New("invalid value")

	// ErrMissingRequired is returned when a required parameter is missing.
	ErrMissingRequired = errors.New("missing required parameter")

	// ErrInvalidFormat is returned when a parameter has an invalid format.
	ErrInvalidFormat = errors.New("invalid format")

	// ErrOutOfRange is returned when a parameter value is out of allowed range.
	ErrOutOfRange = errors.New("value out of range")
)

// ValidationError represents a validation error with context about which field failed.
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

// Error returns the formatted error message.
func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation error for field %q: %s (%v)", e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("validation error for field %q: %s", e.Field, e.Message)
}

// Unwrap returns the wrapped error.
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string, err error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Err:     err,
	}
}

// ValidateInSlice checks if a value is in the allowed slice of values.
func ValidateInSlice(value string, allowed []string, fieldName string) error {
	if value == "" {
		return nil
	}
	for _, a := range allowed {
		if strings.EqualFold(value, a) {
			return nil
		}
	}
	return NewValidationError(fieldName, fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")), ErrInvalidValue)
}

// ValidateSliceValues checks if all values in a slice are in the allowed values.
func ValidateSliceValues(values []string, allowed []string, fieldName string) error {
	for _, value := range values {
		if err := ValidateInSlice(value, allowed, fieldName); err != nil {
			return err
		}
	}
	return nil
}

// ValidateMapKeys checks if all keys in a map are in the allowed keys.
func ValidateMapKeys(m map[string]string, allowedKeys []string, fieldName string) error {
	for key := range m {
		found := false
		for _, allowed := range allowedKeys {
			if strings.EqualFold(key, allowed) {
				found = true
				break
			}
		}
		if !found {
			return NewValidationError(fieldName, fmt.Sprintf("key %q is not allowed. Must be one of: %s", key, strings.Join(allowedKeys, ", ")), ErrInvalidValue)
		}
	}
	return nil
}

// ValidateMapValues checks if all values in a map are in the allowed values.
func ValidateMapValues(m map[string]string, allowedValues []string, fieldName string) error {
	for key, value := range m {
		found := false
		for _, allowed := range allowedValues {
			if strings.EqualFold(value, allowed) {
				found = true
				break
			}
		}
		if !found {
			return NewValidationError(fieldName, fmt.Sprintf("value %q for key %q is not allowed. Must be one of: %s", value, key, strings.Join(allowedValues, ", ")), ErrInvalidValue)
		}
	}
	return nil
}

// ValidateRequired checks if a required field is not empty.
func ValidateRequired(value, fieldName string) error {
	if value == "" {
		return NewValidationError(fieldName, "is required", ErrMissingRequired)
	}
	return nil
}

// ValidateMaxLength checks if a string exceeds maximum length.
func ValidateMaxLength(value string, maxLen int, fieldName string) error {
	if len(value) > maxLen {
		return NewValidationError(fieldName, fmt.Sprintf("exceeds maximum length of %d characters", maxLen), ErrOutOfRange)
	}
	return nil
}

// ValidateMinLength checks if a string is below minimum length.
func ValidateMinLength(value string, minLen int, fieldName string) error {
	if value != "" && len(value) < minLen {
		return NewValidationError(fieldName, fmt.Sprintf("must be at least %d characters", minLen), ErrOutOfRange)
	}
	return nil
}

// ValidateRange checks if a numeric value is within a range.
func ValidateRange(value, min, max int64, fieldName string) error {
	if value < min || value > max {
		return NewValidationError(fieldName, fmt.Sprintf("must be between %d and %d", min, max), ErrOutOfRange)
	}
	return nil
}

// ValidatePagination validates common pagination parameters.
func ValidatePagination(page, pageSize int64) error {
	if page != 0 && page < 1 {
		return NewValidationError("Page", "must be greater than 0", ErrOutOfRange)
	}
	if pageSize != 0 && (pageSize < 1 || pageSize > 500) {
		return NewValidationError("PageSize", "must be between 1 and 500", ErrOutOfRange)
	}
	return nil
}
