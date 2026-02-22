package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	// keyValueParts is the number of parts expected when splitting a key=value string.
	keyValueParts = 2
)

// TriStateBool implements pflag.Value for a tri-state boolean (*bool).
// Accepts "true", "false", or "" (unset). When the flag is not provided,
// the underlying *bool remains nil.
type TriStateBool struct {
	// target is the pointer to the *bool field in the option struct.
	target **bool
	// isSet tracks whether the flag was explicitly set.
	isSet bool
}

// NewTriStateBool creates a new TriStateBool bound to the given **bool target.
func NewTriStateBool(target **bool) *TriStateBool {
	return &TriStateBool{target: target, isSet: false}
}

// String returns the string representation of the current value.
func (t *TriStateBool) String() string {
	if t.target == nil || *t.target == nil {
		return ""
	}

	return strconv.FormatBool(**t.target)
}

// Set parses and sets the value from a string.
func (t *TriStateBool) Set(val string) error {
	t.isSet = true

	switch strings.ToLower(val) {
	case "true":
		v := true
		*t.target = &v
	case "false":
		v := false
		*t.target = &v
	case "":
		*t.target = nil
	default:
		return fmt.Errorf("invalid boolean value: %q (must be true, false, or empty)", val)
	}

	return nil
}

// Type returns the type name for help text.
func (t *TriStateBool) Type() string {
	return "tristate"
}

// IsSet returns whether the flag was explicitly set by the user.
func (t *TriStateBool) IsSet() bool {
	return t.isSet
}

// StringMapValue implements pflag.Value for map[string]string fields.
// Accepts comma-separated KEY=VALUE pairs (e.g., "MAINTAINABILITY=HIGH,SECURITY=LOW").
type StringMapValue struct {
	// target is the pointer to the map[string]string field.
	target *map[string]string
}

// NewStringMapValue creates a new StringMapValue bound to the given target.
func NewStringMapValue(target *map[string]string) *StringMapValue {
	return &StringMapValue{target: target}
}

// String returns the string representation of the current value.
func (m *StringMapValue) String() string {
	if m.target == nil || *m.target == nil {
		return ""
	}

	pairs := make([]string, 0, len(*m.target))
	for k, v := range *m.target {
		pairs = append(pairs, k+"="+v)
	}

	return strings.Join(pairs, ",")
}

// Set parses and sets the value from a string of comma-separated KEY=VALUE pairs.
func (m *StringMapValue) Set(val string) error {
	if val == "" {
		return nil
	}

	result := make(map[string]string)

	for pair := range strings.SplitSeq(val, ",") {
		parts := strings.SplitN(pair, "=", keyValueParts)
		if len(parts) != keyValueParts {
			return fmt.Errorf("invalid key=value pair: %q (expected KEY=VALUE)", pair)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return fmt.Errorf("empty key in pair: %q", pair)
		}

		result[key] = value
	}

	*m.target = result

	return nil
}

// Type returns the type name for help text.
func (m *StringMapValue) Type() string {
	return "KEY=VAL,..."
}

// JSONMapValue implements pflag.Value for JSONEncodedMap (map[string]any) fields.
// Accepts a raw JSON string (e.g., '{"key":"value","count":42}').
type JSONMapValue struct {
	// target is the pointer to the map[string]any field.
	target *map[string]any
}

// NewJSONMapValue creates a new JSONMapValue bound to the given target.
func NewJSONMapValue(target *map[string]any) *JSONMapValue {
	return &JSONMapValue{target: target}
}

// String returns the string representation of the current value.
func (j *JSONMapValue) String() string {
	if j.target == nil || *j.target == nil {
		return ""
	}

	data, err := json.Marshal(*j.target)
	if err != nil {
		return ""
	}

	return string(data)
}

// Set parses and sets the value from a raw JSON string.
func (j *JSONMapValue) Set(val string) error {
	if val == "" {
		return nil
	}

	result := make(map[string]any)

	err := json.Unmarshal([]byte(val), &result)
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	*j.target = result

	return nil
}

// Type returns the type name for help text.
func (j *JSONMapValue) Type() string {
	return "JSON"
}

// CloseBody closes the body of an http.Response safely.
// If the response or body is nil, it does nothing.
// The error return value of Close is intentionally ignored.
func CloseBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
