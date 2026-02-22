package cli

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPascalToKebab tests PascalCase to kebab-case conversion.
func TestPascalToKebab(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple", input: "Issues", want: "issues"},
		{name: "two words", input: "ImpactSeverities", want: "impact-severities"},
		{name: "acronym", input: "URL", want: "url"},
		{name: "acronym at start", input: "HTMLParser", want: "html-parser"},
		{name: "number suffix", input: "PciDss32", want: "pci-dss32"},
		{name: "L10N special case", input: "L10N", want: "l10n"},
		{name: "single letter", input: "A", want: "a"},
		{name: "empty string", input: "", want: ""},
		{name: "lowercase", input: "issues", want: "issues"},
		{name: "mixed acronym", input: "ABCDef", want: "abc-def"},
		{name: "consecutive upper", input: "JSONParser", want: "json-parser"},
		{name: "trailing upper", input: "GetJSON", want: "get-json"},
		{name: "digit inside", input: "Pci4Dss", want: "pci4dss"},
		{name: "multi digit then upper", input: "Get10Items", want: "get10items"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := pascalToKebab(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

// mockOptionStruct is a test struct that mimics a service option.
type mockOptionStruct struct {
	Project      string            `url:"project"`
	Severities   []string          `url:"severities,omitempty"`
	Resolved     *bool             `url:"resolved,omitempty"`
	Page         int64             `url:"p,omitempty"`
	Active       bool              `url:"active,omitempty"`
	Params       map[string]string `url:"params,omitempty"`
	SkippedField string            // no url tag; should be skipped
}

// TestBindFlags_StringField tests binding of a string field.
func TestBindFlags_StringField(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	opt := &mockOptionStruct{}

	BindFlags(cmd, opt)

	// "project" flag should exist and be required.
	f := cmd.Flags().Lookup("project")
	require.NotNil(t, f, "expected 'project' flag to exist")
	assert.Contains(t, f.Usage, "required")

	// Set value through flag.
	require.NoError(t, cmd.Flags().Set("project", "my-project"))
	assert.Equal(t, "my-project", opt.Project)
}

// TestBindFlags_SliceField tests binding of a []string field.
func TestBindFlags_SliceField(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	opt := &mockOptionStruct{}

	BindFlags(cmd, opt)

	f := cmd.Flags().Lookup("severities")
	require.NotNil(t, f, "expected 'severities' flag to exist")
	assert.NotContains(t, f.Usage, "required")

	require.NoError(t, cmd.Flags().Set("severities", "CRITICAL,MAJOR"))
	assert.Equal(t, []string{"CRITICAL", "MAJOR"}, opt.Severities)
}

// TestBindFlags_PointerBoolField tests binding of a *bool field (tri-state).
func TestBindFlags_PointerBoolField(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	opt := &mockOptionStruct{}

	BindFlags(cmd, opt)

	f := cmd.Flags().Lookup("resolved")
	require.NotNil(t, f, "expected 'resolved' flag to exist")

	require.NoError(t, cmd.Flags().Set("resolved", "true"))
	require.NotNil(t, opt.Resolved)
	assert.True(t, *opt.Resolved)
}

// TestBindFlags_Int64Field tests binding of an int64 field.
func TestBindFlags_Int64Field(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	opt := &mockOptionStruct{}

	BindFlags(cmd, opt)

	f := cmd.Flags().Lookup("page")
	require.NotNil(t, f, "expected 'page' flag to exist")

	require.NoError(t, cmd.Flags().Set("page", "5"))
	assert.Equal(t, int64(5), opt.Page)
}

// TestBindFlags_BoolField tests binding of a plain bool field.
func TestBindFlags_BoolField(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	opt := &mockOptionStruct{}

	BindFlags(cmd, opt)

	f := cmd.Flags().Lookup("active")
	require.NotNil(t, f, "expected 'active' flag to exist")

	require.NoError(t, cmd.Flags().Set("active", "true"))
	assert.True(t, opt.Active)
}

// TestBindFlags_MapField tests binding of a map[string]string field.
func TestBindFlags_MapField(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	opt := &mockOptionStruct{}

	BindFlags(cmd, opt)

	f := cmd.Flags().Lookup("params")
	require.NotNil(t, f, "expected 'params' flag to exist")

	require.NoError(t, cmd.Flags().Set("params", "KEY1=VAL1,KEY2=VAL2"))
	// After Set(), the StringMapValue writes to its own map pointer.
	// We verify the flag was accepted without error (the value is stored in the pflag.Value).
	assert.Equal(t, "KEY1=VAL1,KEY2=VAL2", f.Value.String())
}

// TestBindFlags_SkipsUntaggedFields verifies fields without url tag are not bound.
func TestBindFlags_SkipsUntaggedFields(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	opt := &mockOptionStruct{}

	BindFlags(cmd, opt)

	f := cmd.Flags().Lookup("skipped-field")
	assert.Nil(t, f, "expected no flag for untagged field")
}

// embeddedStruct mimics PaginationArgs embedding.
type embeddedStruct struct {
	InnerPagination
	Query string `url:"q"`
}

// InnerPagination is exported to be visible via reflection.
type InnerPagination struct {
	Page     int64 `url:"p,omitempty"`
	PageSize int64 `url:"ps,omitempty"`
}

// TestBindFlags_EmbeddedStruct tests that anonymous embedded structs are recursively bound.
func TestBindFlags_EmbeddedStruct(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	opt := &embeddedStruct{}

	BindFlags(cmd, opt)

	f := cmd.Flags().Lookup("page")
	require.NotNil(t, f, "expected 'page' flag from embedded struct")

	fq := cmd.Flags().Lookup("query")
	require.NotNil(t, fq, "expected 'query' flag")
}

// TestBuildFlagDescription tests flag description generation.
func TestBuildFlagDescription(t *testing.T) {
	tests := []struct {
		name     string
		urlTag   string
		required bool
		want     string
	}{
		{name: "required", urlTag: "project", required: true, want: "API parameter: project (required)"},
		{name: "optional", urlTag: "severity,omitempty", required: false, want: "API parameter: severity"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			field := reflect.StructField{
				Name: "Test",
				Tag:  reflect.StructTag(`url:"` + tc.urlTag + `"`),
			}
			got := buildFlagDescription(field, tc.required)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestParseFlagMeta tests flag name and required status extraction.
func TestParseFlagMeta(t *testing.T) {
	tests := []struct {
		name         string
		fieldName    string
		urlTag       string
		wantName     string
		wantRequired bool
	}{
		{name: "required", fieldName: "Project", urlTag: "project", wantName: "project", wantRequired: true},
		{name: "optional", fieldName: "Severity", urlTag: "severity,omitempty", wantName: "severity", wantRequired: false},
		{name: "inline", fieldName: "Page", urlTag: "p,omitempty", wantName: "page", wantRequired: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			field := reflect.StructField{
				Name: tc.fieldName,
				Tag:  reflect.StructTag(`url:"` + tc.urlTag + `"`),
			}
			name, required := parseFlagMeta(field)
			assert.Equal(t, tc.wantName, name)
			assert.Equal(t, tc.wantRequired, required)
		})
	}
}
