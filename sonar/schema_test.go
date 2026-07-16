package sonar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type schemaTestLeaf struct {
	Name string `json:"name,omitempty"`
}

type schemaTestEmbedded struct {
	Extra string `json:"extra,omitempty"`
}

type schemaTestParent struct {
	schemaTestEmbedded
	Key      string                      `json:"key,omitempty"`
	Ignored  string                      `json:"-"`
	Leaf     schemaTestLeaf              `json:"leaf,omitzero"`
	Leaves   []schemaTestLeaf            `json:"leaves,omitempty"`
	ByName   map[string]schemaTestLeaf   `json:"byName,omitempty"`
	Freeform map[string]any              `json:"freeform,omitempty"`
	Nested   *schemaTestLeaf             `json:"nested,omitempty"`
	Deep     [][]schemaTestLeaf          `json:"deep,omitempty"`
	Various  map[string][]schemaTestLeaf `json:"various,omitempty"`
}

func TestCheckSchema_NoMismatches(t *testing.T) {
	data := []byte(`{"key":"k","extra":"e","leaf":{"name":"n"},"leaves":[{"name":"a"},{"name":"b"}]}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)
	assert.Empty(t, mismatches)
}

func TestCheckSchema_TopLevelUnknownField(t *testing.T) {
	data := []byte(`{"key":"k","ghost":"boo"}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)
	require.Len(t, mismatches, 1)
	assert.Equal(t, "ghost", mismatches[0].Path)
	assert.Equal(t, "GET /x", mismatches[0].Endpoint)
	assert.Contains(t, mismatches[0].GoType, "schemaTestParent")
}

func TestCheckSchema_NestedStructUnknownField(t *testing.T) {
	data := []byte(`{"leaf":{"name":"n","surname":"boo"}}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)
	require.Len(t, mismatches, 1)
	assert.Equal(t, "leaf.surname", mismatches[0].Path)
}

func TestCheckSchema_SliceOfStructsUnknownField(t *testing.T) {
	data := []byte(`{"leaves":[{"name":"a"},{"name":"b","surname":"boo"}]}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)
	require.Len(t, mismatches, 1)
	assert.Equal(t, "leaves[].surname", mismatches[0].Path)
}

func TestCheckSchema_MapOfStructsChecksValues(t *testing.T) {
	data := []byte(`{"byName":{"anyKeyAllowed":{"name":"n","surname":"boo"}}}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)
	require.Len(t, mismatches, 1)
	assert.Equal(t, "byName.anyKeyAllowed.surname", mismatches[0].Path)
}

func TestCheckSchema_FreeformMapAllowsArbitraryKeys(t *testing.T) {
	data := []byte(`{"freeform":{"whatever":"goes","nested":{"a":1}}}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)
	assert.Empty(t, mismatches)
}

func TestCheckSchema_PointerFieldChecked(t *testing.T) {
	data := []byte(`{"nested":{"name":"n","surname":"boo"}}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)
	require.Len(t, mismatches, 1)
	assert.Equal(t, "nested.surname", mismatches[0].Path)
}

func TestCheckSchema_NestedSliceAndMapCombinations(t *testing.T) {
	data := []byte(`{
		"deep": [[{"name":"a","surname":"boo"}]],
		"various": {"k": [{"name":"a"}, {"name":"b","surname":"boo"}]}
	}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)

	paths := make([]string, 0, len(mismatches))
	for _, m := range mismatches {
		paths = append(paths, m.Path)
	}

	assert.ElementsMatch(t, []string{"deep[][].surname", "various.k[].surname"}, paths)
}

func TestCheckSchema_EmbeddedFieldFlattened(t *testing.T) {
	data := []byte(`{"extra":"e"}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)
	assert.Empty(t, mismatches)
}

func TestCheckSchema_IgnoredFieldNeverMatchesOrFlags(t *testing.T) {
	data := []byte(`{"Ignored":"should not match the Go field, should be flagged instead"}`)

	mismatches, err := CheckSchema("GET /x", data, &schemaTestParent{})
	require.NoError(t, err)
	require.Len(t, mismatches, 1)
	assert.Equal(t, "Ignored", mismatches[0].Path)
}

func TestCheckSchema_NilDestYieldsNoMismatches(t *testing.T) {
	mismatches, err := CheckSchema("GET /x", []byte(`{"anything":"goes"}`), nil)
	require.NoError(t, err)
	assert.Nil(t, mismatches)
}

func TestCheckSchema_NonStructDestYieldsNoMismatches(t *testing.T) {
	var dest string

	mismatches, err := CheckSchema("GET /x", []byte(`{"anything":"goes"}`), &dest)
	require.NoError(t, err)
	assert.Empty(t, mismatches)
}

func TestCheckSchema_InvalidJSONReturnsError(t *testing.T) {
	mismatches, err := CheckSchema("GET /x", []byte(`{not valid json`), &schemaTestParent{})
	require.Error(t, err)
	assert.Nil(t, mismatches)
}

func TestSchemaMismatch_String(t *testing.T) {
	mismatch := SchemaMismatch{Endpoint: "GET /x", GoType: "sonar.Foo", Path: "bar"}
	assert.Equal(t, `GET /x: field "bar" has no match in sonar.Foo`, mismatch.String())
}
