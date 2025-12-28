package strcase

import (
	"testing"
)

func TestToCamel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"single word lowercase", "hello", "Hello"},
		{"single word uppercase", "HELLO", "HELLO"},
		{"snake_case", "hello_world", "HelloWorld"},
		{"kebab-case", "hello-world", "HelloWorld"},
		{"space separated", "hello world", "HelloWorld"},
		{"mixed case", "hello_World", "HelloWorld"},
		{"with numbers", "hello_123_world", "Hello123World"},
		{"multiple underscores", "hello__world", "HelloWorld"},
		{"leading underscore", "_hello_world", "HelloWorld"},
		{"trailing underscore", "hello_world_", "HelloWorld"},
		{"all uppercase with underscore", "HELLO_WORLD", "HELLOWORLD"},
		{"already camel case", "HelloWorld", "HelloWorld"},
		{"api endpoint", "api/projects", "Apiprojects"}, // slash is not treated as separator
		{"numbers at start", "123_hello", "123Hello"},
		{"complex path", "get_user_by_id", "GetUserById"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCamel(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamel(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToLowerCamel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"single word lowercase", "hello", "hello"},
		{"single word uppercase", "HELLO", "hELLO"},
		{"snake_case", "hello_world", "helloWorld"},
		{"kebab-case", "hello-world", "helloWorld"},
		{"space separated", "hello world", "helloWorld"},
		{"already lower camel", "helloWorld", "helloWorld"},
		{"with numbers", "hello_123_world", "hello123World"},
		{"uppercase start", "Hello_world", "helloWorld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLowerCamel(tt.input)
			if result != tt.expected {
				t.Errorf("ToLowerCamel(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAddWordBoundariesToNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no numbers", "hello", "hello"},
		{"number at end", "hello123", "hello 123 "}, // regex adds trailing space due to optional 3rd group
		{"number in middle", "hello123world", "hello 123 world"},
		{"multiple numbers", "a1b2c3d", "a 1 b2c 3 d"}, // 2 is consumed by b's match
		{"consecutive numbers", "abc123def", "abc 123 def"},
		{"only numbers", "123", "123"},
		{"number at start", "123abc", "123abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addWordBoundariesToNumbers(tt.input)
			if result != tt.expected {
				t.Errorf("addWordBoundariesToNumbers(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkToCamel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToCamel("hello_world_foo_bar_baz")
	}
}

func BenchmarkToLowerCamel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToLowerCamel("hello_world_foo_bar_baz")
	}
}
