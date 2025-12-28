package validation

import (
	"strings"
	"testing"
)

func TestNewFieldRequired(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected Error
	}{
		{
			name:  "simple field",
			field: "username",
			expected: Error{
				Type:  ErrorTypeRequired,
				Field: "username",
			},
		},
		{
			name:  "nested field",
			field: "user.email",
			expected: Error{
				Type:  ErrorTypeRequired,
				Field: "user.email",
			},
		},
		{
			name:  "empty field",
			field: "",
			expected: Error{
				Type:  ErrorTypeRequired,
				Field: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewFieldRequired(tt.field)
			if err.Type != tt.expected.Type {
				t.Errorf("NewFieldRequired(%q).Type = %q, want %q", tt.field, err.Type, tt.expected.Type)
			}
			if err.Field != tt.expected.Field {
				t.Errorf("NewFieldRequired(%q).Field = %q, want %q", tt.field, err.Field, tt.expected.Field)
			}
		})
	}
}

func TestNewFieldInvalidValue(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected Error
	}{
		{
			name:  "simple field",
			field: "age",
			expected: Error{
				Type:  ErrorInvalidValue,
				Field: "age",
			},
		},
		{
			name:  "nested field",
			field: "config.timeout",
			expected: Error{
				Type:  ErrorInvalidValue,
				Field: "config.timeout",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewFieldInvalidValue(tt.field)
			if err.Type != tt.expected.Type {
				t.Errorf("NewFieldInvalidValue(%q).Type = %q, want %q", tt.field, err.Type, tt.expected.Type)
			}
			if err.Field != tt.expected.Field {
				t.Errorf("NewFieldInvalidValue(%q).Field = %q, want %q", tt.field, err.Field, tt.expected.Field)
			}
		})
	}
}

func TestNewFieldInvalidValueWithReason(t *testing.T) {
	tests := []struct {
		name           string
		field          string
		reason         string
		expectedType   ErrorType
		expectedField  string
		expectedReason string
	}{
		{
			name:           "with reason",
			field:          "port",
			reason:         "must be between 1 and 65535",
			expectedType:   ErrorInvalidValue,
			expectedField:  "port",
			expectedReason: "must be between 1 and 65535",
		},
		{
			name:           "empty reason",
			field:          "url",
			reason:         "",
			expectedType:   ErrorInvalidValue,
			expectedField:  "url",
			expectedReason: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewFieldInvalidValueWithReason(tt.field, tt.reason)
			if err.Type != tt.expectedType {
				t.Errorf("Type = %q, want %q", err.Type, tt.expectedType)
			}
			if err.Field != tt.expectedField {
				t.Errorf("Field = %q, want %q", err.Field, tt.expectedField)
			}
			if err.Reason != tt.expectedReason {
				t.Errorf("Reason = %q, want %q", err.Reason, tt.expectedReason)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      Error
		contains []string
	}{
		{
			name: "required error",
			err: Error{
				Type:  ErrorTypeRequired,
				Field: "username",
			},
			contains: []string{"Required", "username"},
		},
		{
			name: "invalid value error",
			err: Error{
				Type:  ErrorInvalidValue,
				Field: "age",
			},
			contains: []string{"Invalid", "age"},
		},
		{
			name: "invalid value with reason",
			err: Error{
				Type:   ErrorInvalidValue,
				Field:  "port",
				Reason: "must be positive",
			},
			contains: []string{"Invalid", "port", "must be positive"},
		},
		{
			name: "unknown type",
			err: Error{
				Type:  "UnknownType",
				Field: "field",
			},
			contains: []string{"UnknownType", "field"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			for _, substr := range tt.contains {
				if !strings.Contains(msg, substr) {
					t.Errorf("Error() = %q, should contain %q", msg, substr)
				}
			}
		})
	}
}

func TestErrorType_Constants(t *testing.T) {
	if ErrorTypeRequired != "FieldValueRequired" {
		t.Errorf("ErrorTypeRequired = %q, want %q", ErrorTypeRequired, "FieldValueRequired")
	}
	if ErrorInvalidValue != "InvalidValue" {
		t.Errorf("ErrorInvalidValue = %q, want %q", ErrorInvalidValue, "InvalidValue")
	}
}
