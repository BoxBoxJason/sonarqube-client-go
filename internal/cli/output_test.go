package cli

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sampleStruct is a test struct for output formatting.
type sampleStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// TestFormatOutput_Nil tests that nil values produce no output.
func TestFormatOutput_Nil(t *testing.T) {
	var buf bytes.Buffer

	err := FormatOutput(&buf, nil, OutputJSON)
	require.NoError(t, err)
	assert.Empty(t, buf.String())
}

// TestFormatOutput_RawBytes tests that raw byte slices are written directly.
func TestFormatOutput_RawBytes(t *testing.T) {
	var buf bytes.Buffer
	data := []byte("raw binary content")

	err := FormatOutput(&buf, data, OutputJSON)
	require.NoError(t, err)
	assert.Equal(t, "raw binary content", buf.String())
}

// TestFormatOutput_RawString tests that *string values are printed as lines.
func TestFormatOutput_RawString(t *testing.T) {
	var buf bytes.Buffer
	val := "hello world"

	err := FormatOutput(&buf, &val, OutputJSON)
	require.NoError(t, err)
	assert.Equal(t, "hello world\n", buf.String())
}

// TestFormatOutput_JSON tests JSON output of a struct.
func TestFormatOutput_JSON(t *testing.T) {
	var buf bytes.Buffer
	data := &sampleStruct{Name: "test", Value: 42}

	err := FormatOutput(&buf, data, OutputJSON)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), `"name": "test"`)
	assert.Contains(t, buf.String(), `"value": 42`)
}

// TestFormatOutput_YAML tests YAML output of a struct.
func TestFormatOutput_YAML(t *testing.T) {
	var buf bytes.Buffer
	data := &sampleStruct{Name: "test", Value: 42}

	err := FormatOutput(&buf, data, OutputYAML)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "name: test")
	assert.Contains(t, buf.String(), "value: 42")
}

// TestFormatOutput_Table_Struct tests table output for a single struct.
func TestFormatOutput_Table_Struct(t *testing.T) {
	var buf bytes.Buffer
	data := &sampleStruct{Name: "test", Value: 42}

	err := FormatOutput(&buf, data, OutputTable)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "FIELD")
	assert.Contains(t, output, "VALUE")
	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "test")
}

// TestFormatOutput_Table_Slice tests table output for a slice of structs.
func TestFormatOutput_Table_Slice(t *testing.T) {
	var buf bytes.Buffer
	data := []sampleStruct{
		{Name: "alpha", Value: 1},
		{Name: "beta", Value: 2},
	}

	err := FormatOutput(&buf, data, OutputTable)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, "VALUE")
	assert.Contains(t, output, "alpha")
	assert.Contains(t, output, "beta")
}

// TestFormatOutput_Table_EmptySlice tests table output for an empty slice.
func TestFormatOutput_Table_EmptySlice(t *testing.T) {
	var buf bytes.Buffer
	data := []sampleStruct{}

	err := FormatOutput(&buf, data, OutputTable)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "(no results)")
}

// TestFormatOutput_DefaultFormat tests that an unknown format falls back to JSON.
func TestFormatOutput_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	data := &sampleStruct{Name: "test", Value: 42}

	err := FormatOutput(&buf, data, "unknown")
	require.NoError(t, err)
	assert.Contains(t, buf.String(), `"name": "test"`)
}

// TestExtractHeaders tests header extraction from struct types.
func TestExtractHeaders(t *testing.T) {
	type withJSON struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	headers := extractHeaders(reflect.TypeOf(withJSON{}))
	assert.Equal(t, []string{"ID", "NAME"}, headers)
}

// TestExtractHeaders_NoJSONTag tests header extraction without json tags.
func TestExtractHeaders_NoJSONTag(t *testing.T) {
	type noTag struct {
		First  string
		Second int
	}

	headers := extractHeaders(reflect.TypeOf(noTag{}))
	assert.Equal(t, []string{"FIRST", "SECOND"}, headers)
}

// TestRenderTable tests the ASCII table renderer.
func TestRenderTable(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"A", "B"}
	rows := [][]string{
		{"hello", "world"},
		{"foo", "bar"},
	}

	renderTable(&buf, headers, rows)

	output := buf.String()
	assert.Contains(t, output, "A")
	assert.Contains(t, output, "B")
	assert.Contains(t, output, "hello")
	assert.Contains(t, output, "bar")
	assert.Contains(t, output, "---")
}

// TestRenderTable_Empty tests that an empty headers slice produces no output.
func TestRenderTable_Empty(t *testing.T) {
	var buf bytes.Buffer

	renderTable(&buf, nil, nil)
	assert.Empty(t, buf.String())
}
