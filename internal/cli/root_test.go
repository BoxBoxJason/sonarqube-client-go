package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildRootCommand verifies root command configuration.
func TestBuildRootCommand(t *testing.T) {
	flags := &globalFlags{}
	cmd := buildRootCommand(flags)

	assert.Equal(t, "sonar-cli", cmd.Use)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Verify global flags are registered.
	globalFlagNames := []string{"url", "token", "username", "password", "output", "timeout"}
	for _, name := range globalFlagNames {
		f := cmd.PersistentFlags().Lookup(name)
		assert.NotNil(t, f, "expected persistent flag %q", name)
	}
}

// TestGlobalFlags_Defaults tests default values for global flags.
func TestGlobalFlags_Defaults(t *testing.T) {
	flags := &globalFlags{}
	cmd := buildRootCommand(flags)

	// Check URL flag has no default (must be provided via flag or env var).
	urlFlag := cmd.PersistentFlags().Lookup("url")
	require.NotNil(t, urlFlag)
	assert.Equal(t, "", urlFlag.DefValue)

	timeoutFlag := cmd.PersistentFlags().Lookup("timeout")
	require.NotNil(t, timeoutFlag)
	assert.Equal(t, defaultTimeout.String(), timeoutFlag.DefValue)
}

// TestOutputFormatFlag tests the custom output format flag validation.
func TestOutputFormatFlag(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    OutputFormat
	}{
		{name: "json", input: "json", want: OutputJSON},
		{name: "table", input: "table", want: OutputTable},
		{name: "yaml", input: "yaml", want: OutputYAML},
		{name: "invalid", input: "xml", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			target := defaultOutputFormat
			f := &outputFormatFlag{target: &target}

			err := f.Set(tc.input)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, target)
		})
	}
}

// TestOutputFormatFlag_String tests string representation.
func TestOutputFormatFlag_String(t *testing.T) {
	target := OutputTable
	f := &outputFormatFlag{target: &target}

	assert.Equal(t, "table", f.String())
}

// TestOutputFormatFlag_String_Nil tests nil target.
func TestOutputFormatFlag_String_Nil(t *testing.T) {
	f := &outputFormatFlag{target: nil}
	assert.Equal(t, string(defaultOutputFormat), f.String())
}

// TestOutputFormatFlag_Type tests the type name.
func TestOutputFormatFlag_Type(t *testing.T) {
	f := &outputFormatFlag{}
	assert.Equal(t, "format", f.Type())
}

// TestClientFromContext_NoContext tests error when no client in context.
func TestClientFromContext_NoContext(t *testing.T) {
	flags := &globalFlags{}
	cmd := buildRootCommand(flags)

	// Set a context so cmd.Context() doesn't return nil.
	cmd.SetContext(context.Background())

	_, err := clientFromContext(cmd)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}
