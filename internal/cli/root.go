package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/boxboxjason/sonarqube-client-go/sonar"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	// defaultTimeout is the default HTTP request timeout.
	defaultTimeout = 30 * time.Second
	// defaultOutputFormat is the default output format for CLI responses.
	defaultOutputFormat = OutputJSON
	// completeDirective is the special argument for shell completion.
	completeDirective = "__complete"
	// completeNoDescDirective is the special argument for shell completion without descriptions.
	completeNoDescDirective = "__completeNoDesc"
)

// contextKey is a typed key for storing values in cobra command context.
type contextKey string

const clientContextKey contextKey = "sonar-client"

// globalFlags holds the global CLI flag values.
type globalFlags struct {
	url      string
	token    string
	username string
	password string
	output   OutputFormat
	timeout  time.Duration
}

// Execute creates the root command, registers all subcommands, and runs the CLI.
func Execute() error {
	// Initialize logger
	initLogger()

	defer Sync() //nolint:errcheck

	flags := &globalFlags{ //nolint:exhaustruct // fields are set by Cobra flag binding
		output: defaultOutputFormat,
	}

	rootCmd := buildRootCommand(flags)

	RegisterAllCommands(rootCmd, &flags.output)

	return rootCmd.Execute() //nolint:wrapcheck // errors are logged at source (initClient, runMethodCommand)
}

// buildRootCommand creates and configures the root Cobra command with global flags.
func buildRootCommand(flags *globalFlags) *cobra.Command {
	rootCmd := &cobra.Command{ //nolint:exhaustruct // only setting fields relevant to root command
		Use:   "sonar-cli",
		Short: "CLI for the SonarQube API",
		Long: `sonar-cli is a command-line interface for the SonarQube API.
It wraps the sonarqube-client-go SDK, providing access to all SonarQube
API endpoints from the terminal.

Usage:
  sonar-cli [global flags] <service> <method> [flags]

Examples:
  sonar-cli --token mytoken issues search --severities CRITICAL,MAJOR
  sonar-cli --url http://sonar:9000 --token mytoken projects search --all
  sonar-cli --output table qualitygates list`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initClient(cmd, args, flags)
		},
		Version: versionInfo(),
	}

	registerGlobalFlags(rootCmd, flags)

	// Enable shell completion validation for all shells
	rootCmd.CompletionOptions.DisableDescriptions = false
	rootCmd.CompletionOptions.HiddenDefaultCmd = false

	return rootCmd
}

// registerGlobalFlags adds persistent (global) flags to the root command.
func registerGlobalFlags(cmd *cobra.Command, flags *globalFlags) {
	persistentFlags := cmd.PersistentFlags()

	persistentFlags.StringVar(&flags.url, "url", os.Getenv("SONAR_CLI_URL"), "SonarQube server URL (also read from SONAR_CLI_URL env var)")
	persistentFlags.StringVar(&flags.token, "token", os.Getenv("SONAR_CLI_TOKEN"), "Authentication token (also read from SONAR_CLI_TOKEN env var)")
	persistentFlags.StringVar(&flags.username, "username", os.Getenv("SONAR_CLI_USERNAME"), "Username for basic authentication (also read from SONAR_CLI_USERNAME env var)")
	persistentFlags.StringVar(&flags.password, "password", os.Getenv("SONAR_CLI_PASSWORD"), "Password for basic authentication (also read from SONAR_CLI_PASSWORD env var)")
	persistentFlags.DurationVar(&flags.timeout, "timeout", defaultTimeout, "HTTP request timeout")

	// Output format with custom validation.
	flags.output = defaultOutputFormat

	persistentFlags.Var(&outputFormatFlag{target: &flags.output}, "output", "Output format: json, table, yaml")
}

// initClient creates a sonar.Client from global flags and stores it in the command context.
// This runs as PersistentPreRunE so the client is available to all subcommands.
// It skips initialization for the completion command and shell completion directives.
func initClient(cmd *cobra.Command, args []string, globalFlags *globalFlags) error {
	if shouldSkipClientInit(cmd, args) {
		return nil
	}

	opts := &sonar.ClientCreateOption{} //nolint:exhaustruct // fields set conditionally below

	if globalFlags.url == "" {
		err := errors.New("server URL must be provided via --url flag or SONAR_CLI_URL env var")
		Logger().Error("missing required configuration", zap.Error(err))

		return err
	}

	opts.URL = &globalFlags.url
	setClientOptionalFields(opts, globalFlags)

	httpClient := &http.Client{ //nolint:exhaustruct // only Timeout is needed
		Timeout: globalFlags.timeout,
	}
	opts.HttpClient = httpClient

	client, err := sonar.NewClient(opts)
	if err != nil {
		Logger().Error("failed to initialize SonarQube client", zap.Error(err))

		return fmt.Errorf("failed to create SonarQube client: %w", err)
	}

	cmd.SetContext(context.WithValue(cmd.Context(), clientContextKey, client))

	return nil
}

// shouldSkipClientInit checks if client initialization should be skipped.
func shouldSkipClientInit(cmd *cobra.Command, args []string) bool {
	return isCompletionCommand(cmd) || isCompletionDirective(args)
}

// isCompletionCommand checks if the command is a completion or __complete command.
func isCompletionCommand(cmd *cobra.Command) bool {
	if cmd.Name() == "completion" || cmd.Name() == completeDirective || cmd.Name() == completeNoDescDirective {
		return true
	}

	if cmd.Parent() != nil {
		parentName := cmd.Parent().Name()
		if parentName == "completion" || parentName == completeDirective || parentName == completeNoDescDirective {
			return true
		}
	}

	return false
}

// isCompletionDirective checks if the first argument is a completion directive.
func isCompletionDirective(args []string) bool {
	if len(args) > 0 && (args[0] == completeDirective || args[0] == completeNoDescDirective) {
		return true
	}

	return false
}

// setClientOptionalFields sets optional authentication fields on the client options.
func setClientOptionalFields(opts *sonar.ClientCreateOption, globalFlags *globalFlags) {
	if globalFlags.token != "" {
		opts.Token = &globalFlags.token
	}

	if globalFlags.username != "" {
		opts.Username = &globalFlags.username
	}

	if globalFlags.password != "" {
		opts.Password = &globalFlags.password
	}
}

// clientFromContext retrieves the sonar.Client from the command's context.
func clientFromContext(cmd *cobra.Command) (*sonar.Client, error) {
	val := cmd.Context().Value(clientContextKey)
	if val == nil {
		return nil, errors.New("sonarqube client not initialized (this is a bug)")
	}

	client, ok := val.(*sonar.Client)
	if !ok {
		return nil, errors.New("invalid client type in context (this is a bug)")
	}

	return client, nil
}

// outputFormatFlag implements pflag.Value for the --output flag with validation.
type outputFormatFlag struct {
	target *OutputFormat
}

// String returns the current output format as a string.
func (f *outputFormatFlag) String() string {
	if f.target == nil {
		return string(defaultOutputFormat)
	}

	return string(*f.target)
}

// Set validates and sets the output format.
func (f *outputFormatFlag) Set(val string) error {
	switch OutputFormat(val) {
	case OutputJSON, OutputTable, OutputYAML:
		*f.target = OutputFormat(val)

		return nil
	default:
		return fmt.Errorf("invalid output format %q: must be one of json, table, yaml", val)
	}
}

// Type returns the type name for help text.
func (f *outputFormatFlag) Type() string {
	return "format"
}
