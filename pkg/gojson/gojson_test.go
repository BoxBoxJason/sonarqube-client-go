package gojson

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseJson(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid simple object",
			input:       `{"name": "test", "value": 123}`,
			expectError: false,
		},
		{
			name:        "valid array",
			input:       `[{"name": "item1"}, {"name": "item2"}]`,
			expectError: false,
		},
		{
			name:        "valid nested object",
			input:       `{"user": {"name": "john", "age": 30}}`,
			expectError: false,
		},
		{
			name:        "empty object",
			input:       `{}`,
			expectError: false,
		},
		{
			name:        "invalid json",
			input:       `{"name": "test"`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.input))
			result, err := ParseJson(reader)

			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && result == nil {
				t.Error("expected result but got nil")
			}
		})
	}
}

func TestParseYaml(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name: "valid yaml",
			input: `name: test
value: 123`,
			expectError: false,
		},
		{
			name: "valid nested yaml",
			input: `user:
  name: john
  age: 30`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.input))
			result, err := ParseYaml(reader)

			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && result == nil {
				t.Error("expected result but got nil")
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		structName    string
		tags          []string
		subStruct     bool
		convertFloats bool
		wantContains  []string
		wantError     bool
	}{
		{
			name:          "simple struct",
			input:         `{"name": "test", "count": 42}`,
			structName:    "MyStruct",
			tags:          []string{"json"},
			subStruct:     false,
			convertFloats: true,
			wantContains:  []string{"type MyStruct struct", "Name", "Count", "json:"},
			wantError:     false,
		},
		{
			name:          "nested struct without substruct",
			input:         `{"user": {"name": "john", "age": 30}}`,
			structName:    "Response",
			tags:          []string{"json"},
			subStruct:     false,
			convertFloats: true,
			wantContains:  []string{"type Response struct", "User"},
			wantError:     false,
		},
		{
			name:          "with multiple tags",
			input:         `{"id": 1}`,
			structName:    "Item",
			tags:          []string{"json", "xml"},
			subStruct:     false,
			convertFloats: true,
			wantContains:  []string{"json:", "xml:"},
			wantError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.input))
			result, err := Generate(reader, ParseJson, tt.structName, tt.tags, tt.subStruct, tt.convertFloats)

			if tt.wantError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.wantError {
				resultStr := string(result)
				for _, want := range tt.wantContains {
					if !strings.Contains(resultStr, want) {
						t.Errorf("result should contain %q, got: %s", want, resultStr)
					}
				}
			}
		})
	}
}

func TestFmtFieldName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_name", "UserName"},
		{"userId", "UserID"},
		{"api_key", "APIKey"},
		{"simple", "Simple"},
		{"URL", "URL"},
		{"http_url", "HTTPURL"},
		{"id", "ID"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := FmtFieldName(tt.input)
			if result != tt.expected {
				t.Errorf("FmtFieldName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStringifyFirstChar(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123field", "one_23field"},
		{"field", "field"},
		{"_field", "_field"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := stringifyFirstChar(tt.input)
			if result != tt.expected {
				t.Errorf("stringifyFirstChar(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDisambiguateFloatInt(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"integer float", float64(42), "int64"},
		{"actual float", float64(3.14), "float64"},
		{"large integer", float64(1000000), "int64"},
		{"negative integer", float64(-42), "int64"},
		{"zero", float64(0), "int64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := disambiguateFloatInt(tt.input)
			if result != tt.expected {
				t.Errorf("disambiguateFloatInt(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMergeElements(t *testing.T) {
	t.Run("merge array of objects", func(t *testing.T) {
		input := []interface{}{
			map[string]interface{}{"name": "a"},
			map[string]interface{}{"name": "b", "age": 30},
		}
		result := mergeElements(input)
		if result == nil {
			t.Error("expected non-nil result")
		}
	})

	t.Run("non-array input", func(t *testing.T) {
		input := map[string]interface{}{"name": "test"}
		result := mergeElements(input)
		if result == nil {
			t.Error("expected non-nil result")
		}
	})
}

func TestConvertKeysToStrings(t *testing.T) {
	input := map[interface{}]interface{}{
		"name": "test",
		"age":  30,
	}

	result := convertKeysToStrings(input)

	if result["name"] != "test" {
		t.Errorf("expected name to be 'test', got %v", result["name"])
	}
	if result["age"] != 30 {
		t.Errorf("expected age to be 30, got %v", result["age"])
	}
}
