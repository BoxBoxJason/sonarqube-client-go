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
	Err     error
	Field   string
	Message string
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string, err error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Err:     err,
	}
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

// IsValueAuthorized checks if a value is in the allowed set of values.
func IsValueAuthorized(value string, allowed map[string]struct{}, fieldName string) error {
	if value == "" {
		return nil
	}

	_, isInSlice := allowed[value]
	if !isInSlice {
		return NewValidationError(fieldName, "must be one of: "+BuildAuthorizedValuesList(allowed), ErrInvalidValue)
	}

	return nil
}

// AreValuesAuthorized checks if all values in a slice are in the allowed values.
func AreValuesAuthorized(values []string, allowed map[string]struct{}, fieldName string) error {
	for _, value := range values {
		_, isInSlice := allowed[value]
		if !isInSlice {
			return NewValidationError(fieldName, fmt.Sprintf("value %q is not allowed. Must be one of: %s", value, BuildAuthorizedValuesList(allowed)), ErrInvalidValue)
		}
	}

	return nil
}

// ValidateMapKeys checks if all keys in a map are in the allowed keys.
func ValidateMapKeys(m map[string]string, allowedKeys map[string]struct{}, fieldName string) error {
	for key := range m {
		_, found := allowedKeys[key]

		if !found {
			return NewValidationError(fieldName, fmt.Sprintf("key %q is not allowed. Must be one of: %s", key, BuildAuthorizedValuesList(allowedKeys)), ErrInvalidValue)
		}
	}

	return nil
}

// ValidateMapValues checks if all values in a map are in the allowed values.
func ValidateMapValues(m map[string]string, allowedValues map[string]struct{}, fieldName string) error {
	for key, value := range m {
		_, found := allowedValues[value]

		if !found {
			return NewValidationError(fieldName, fmt.Sprintf("value %q for key %q is not allowed. Must be one of: %s", value, key, BuildAuthorizedValuesList(allowedValues)), ErrInvalidValue)
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
func ValidateRange(value, minValue, maxValue int64, fieldName string) error {
	if value < minValue || value > maxValue {
		return NewValidationError(fieldName, fmt.Sprintf("must be between %d and %d", minValue, maxValue), ErrOutOfRange)
	}

	return nil
}

// ValidatePagination validates common pagination parameters.
func ValidatePagination(page, pageSize int64) error {
	if page != 0 && page < MinPageSize {
		return NewValidationError("Page", "must be greater than 0", ErrOutOfRange)
	}

	if pageSize != 0 && (pageSize < MinPageSize || pageSize > MaxPageSize) {
		return NewValidationError("PageSize", fmt.Sprintf("must be between %d and %d", MinPageSize, MaxPageSize), ErrOutOfRange)
	}

	return nil
}

// ValidateLanguage checks if the provided language is among the allowed languages.
func ValidateLanguage(language string) error {
	return IsValueAuthorized(language, allowedLanguages, "Language")
}

// ValidateLanguages checks if all provided languages are among the allowed languages.
func ValidateLanguages(languages []string) error {
	return AreValuesAuthorized(languages, allowedLanguages, "Languages")
}

// BuildAuthorizedValuesList builds a comma-separated string of allowed values.
func BuildAuthorizedValuesList(allowed map[string]struct{}) string {
	responseBuilder := strings.Builder{}
	index := 0

	for key := range allowed {
		responseBuilder.WriteString(key)

		if index < len(allowed)-1 {
			responseBuilder.WriteString(", ")
		}

		index++
	}

	return responseBuilder.String()
}
