package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTriStateBool_Set tests parsing of tri-state boolean values.
func TestTriStateBool_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantVal bool
		wantErr bool
	}{
		{name: "true", input: "true", wantNil: false, wantVal: true},
		{name: "false", input: "false", wantNil: false, wantVal: false},
		{name: "TRUE", input: "TRUE", wantNil: false, wantVal: true},
		{name: "FALSE", input: "FALSE", wantNil: false, wantVal: false},
		{name: "empty", input: "", wantNil: true},
		{name: "invalid", input: "maybe", wantErr: true},
		{name: "number", input: "1", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var target *bool
			tsb := NewTriStateBool(&target)

			err := tsb.Set(tc.input)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			if tc.wantNil {
				assert.Nil(t, target)
			} else {
				require.NotNil(t, target)
				assert.Equal(t, tc.wantVal, *target)
			}

			assert.True(t, tsb.IsSet())
		})
	}
}

// TestTriStateBool_String tests string representation.
func TestTriStateBool_String(t *testing.T) {
	var target *bool
	tsb := NewTriStateBool(&target)

	assert.Equal(t, "", tsb.String())

	val := true
	target = &val
	assert.Equal(t, "true", tsb.String())

	val = false
	assert.Equal(t, "false", tsb.String())
}

// TestTriStateBool_Type tests type name.
func TestTriStateBool_Type(t *testing.T) {
	var target *bool
	tsb := NewTriStateBool(&target)

	assert.Equal(t, "tristate", tsb.Type())
}

// TestStringMapValue_Set tests parsing of key=value pairs.
func TestStringMapValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{name: "single pair", input: "KEY=VAL", want: map[string]string{"KEY": "VAL"}},
		{name: "multiple pairs", input: "A=1,B=2,C=3", want: map[string]string{"A": "1", "B": "2", "C": "3"}},
		{name: "empty string", input: "", want: nil},
		{name: "no equals", input: "INVALID", wantErr: true},
		{name: "empty key", input: "=VAL", wantErr: true},
		{name: "value with equals", input: "KEY=A=B", want: map[string]string{"KEY": "A=B"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var target map[string]string
			smv := NewStringMapValue(&target)

			err := smv.Set(tc.input)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, target)
		})
	}
}

// TestStringMapValue_String tests string representation of map values.
func TestStringMapValue_String(t *testing.T) {
	var target map[string]string
	smv := NewStringMapValue(&target)

	assert.Equal(t, "", smv.String())

	target = map[string]string{"A": "1"}
	assert.Equal(t, "A=1", smv.String())
}

// TestStringMapValue_Type tests the type name.
func TestStringMapValue_Type(t *testing.T) {
	var target map[string]string
	smv := NewStringMapValue(&target)

	assert.Equal(t, "KEY=VAL,...", smv.Type())
}

// TestJSONMapValue_Set tests parsing of JSON string values.
func TestJSONMapValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]any
		wantErr bool
	}{
		{name: "valid JSON", input: `{"key":"val"}`, want: map[string]any{"key": "val"}},
		{name: "numeric value", input: `{"count":42}`, want: map[string]any{"count": float64(42)}},
		{name: "empty string", input: "", want: nil},
		{name: "invalid JSON", input: "not json", wantErr: true},
		{name: "array", input: "[1,2,3]", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var target map[string]any
			jmv := NewJSONMapValue(&target)

			err := jmv.Set(tc.input)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, target)
		})
	}
}

// TestJSONMapValue_String tests JSON string representation.
func TestJSONMapValue_String(t *testing.T) {
	var target map[string]any
	jmv := NewJSONMapValue(&target)

	assert.Equal(t, "", jmv.String())

	target = map[string]any{"key": "val"}
	// JSON marshaling produces consistent output for simple maps.
	assert.Equal(t, `{"key":"val"}`, jmv.String())
}

// TestJSONMapValue_Type tests the type name.
func TestJSONMapValue_Type(t *testing.T) {
	var target map[string]any
	jmv := NewJSONMapValue(&target)

	assert.Equal(t, "JSON", jmv.Type())
}
