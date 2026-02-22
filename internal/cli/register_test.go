package cli

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegisterAllCommands verifies that all sonar.Client services are registered as subcommands.
func TestRegisterAllCommands(t *testing.T) {
	format := OutputJSON
	rootCmd := &cobra.Command{Use: "test"}

	RegisterAllCommands(rootCmd, &format)

	// Verify that subcommands were registered (at least a known set).
	commands := rootCmd.Commands()
	assert.NotEmpty(t, commands, "expected subcommands to be registered")

	// Check for some well-known services.
	knownServices := []string{"issues", "projects", "rules", "qualitygates", "system"}
	for _, svc := range knownServices {
		found := false

		for _, cmd := range commands {
			if cmd.Name() == svc {
				found = true

				break
			}
		}

		assert.True(t, found, "expected service %q to be registered", svc)
	}
}

// TestRegisterAllCommands_SubCommands verifies that service commands have method subcommands.
func TestRegisterAllCommands_SubCommands(t *testing.T) {
	format := OutputJSON
	rootCmd := &cobra.Command{Use: "test"}

	RegisterAllCommands(rootCmd, &format)

	// Find the "issues" service command.
	var issuesCmd *cobra.Command

	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "issues" {
			issuesCmd = cmd

			break
		}
	}

	require.NotNil(t, issuesCmd, "expected 'issues' service command")

	// Issues should have at least a "search" subcommand.
	subCommands := issuesCmd.Commands()
	assert.NotEmpty(t, subCommands, "expected issues subcommands")

	found := false

	for _, cmd := range subCommands {
		if cmd.Name() == "search" {
			found = true

			break
		}
	}

	assert.True(t, found, "expected 'search' subcommand on issues")
}

// TestShouldSkipMethod tests method name filtering.
func TestShouldSkipMethod(t *testing.T) {
	tests := []struct {
		name   string
		method string
		skip   bool
	}{
		{name: "normal method", method: "Search", skip: false},
		{name: "validate method", method: "ValidateInput", skip: true},
		{name: "exact validate", method: "Validate", skip: true},
		{name: "create", method: "Create", skip: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldSkipMethod(tc.method)
			assert.Equal(t, tc.skip, got)
		})
	}
}

// TestBuildServiceCommand_EmptyService tests that a service with no valid methods returns nil.
func TestBuildServiceCommand_EmptyService(t *testing.T) {
	// emptyService has no exported methods that match the pattern.
	type emptyService struct{}

	format := OutputJSON
	cmd := buildServiceCommand("Empty", reflect.TypeOf(&emptyService{}), &format)
	assert.Nil(t, cmd, "expected nil command for service with no methods")
}

// TestGetServiceDescription tests service description retrieval.
func TestGetServiceDescription(t *testing.T) {
	// Known service should return a description.
	desc := GetServiceDescription("Issues")
	assert.NotEmpty(t, desc)
	assert.NotEqual(t, "Manages Issues operations.", desc, "should return a custom description, not the fallback")

	// Unknown service should return a fallback.
	desc = GetServiceDescription("UnknownTestService")
	assert.Contains(t, desc, "UnknownTestService")
}

// TestGetMethodDescription tests method description retrieval.
func TestGetMethodDescription(t *testing.T) {
	// Known method should return a description.
	desc := GetMethodDescription("Issues", "Search")
	assert.NotEmpty(t, desc)

	// Unknown method should return a fallback.
	desc = GetMethodDescription("Issues", "NonexistentMethod")
	assert.Contains(t, desc, "NonexistentMethod")
}
