package cli

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/boxboxjason/sonarqube-client-go/sonar"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// streamingMethods lists methods that return streaming responses and need special handling.
// Key format: "ServiceFieldName.MethodName".
//
//nolint:gochecknoglobals // constant configuration set
var streamingMethods = map[string]struct{}{
	"Push.SonarlintEvents": {},
}

// skipMethods lists method name prefixes that should not become CLI commands.
//
//nolint:gochecknoglobals // constant configuration set
var skipMethods = map[string]struct{}{
	"Validate": {},
}

// RegisterAllCommands discovers all services on the sonar.Client and registers
// Cobra subcommands for each public method.
func RegisterAllCommands(rootCmd *cobra.Command, format *OutputFormat) {
	clientType := reflect.TypeFor[sonar.Client]()

	for fieldIdx := range clientType.NumField() {
		field := clientType.Field(fieldIdx)

		// Only look at exported pointer-to-struct fields (the service fields).
		if !field.IsExported() || field.Type.Kind() != reflect.Ptr || field.Type.Elem().Kind() != reflect.Struct {
			continue
		}

		serviceName := field.Name
		serviceType := field.Type

		serviceCmd := buildServiceCommand(serviceName, serviceType, format)
		if serviceCmd != nil {
			rootCmd.AddCommand(serviceCmd)
		}
	}
}

// buildServiceCommand creates a Cobra command for a service, with subcommands for each method.
func buildServiceCommand(serviceName string, serviceType reflect.Type, format *OutputFormat) *cobra.Command {
	kebabName := pascalToKebab(serviceName)
	description := GetServiceDescription(serviceName)

	serviceCmd := &cobra.Command{ //nolint:exhaustruct // only Use/Short/Long are needed
		Use:   kebabName,
		Short: description,
		Long:  fmt.Sprintf("Commands for the %s service.", serviceName),
	}

	methodCount := 0

	for methodIdx := range serviceType.NumMethod() {
		method := serviceType.Method(methodIdx)

		if shouldSkipMethod(method.Name) {
			continue
		}

		methodCmd := buildMethodCommand(serviceName, serviceType, method, format)
		if methodCmd != nil {
			serviceCmd.AddCommand(methodCmd)

			methodCount++
		}
	}

	if methodCount == 0 {
		return nil
	}

	return serviceCmd
}

// shouldSkipMethod returns true if the method should not be registered as a CLI command.
func shouldSkipMethod(name string) bool {
	for prefix := range skipMethods {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

// buildMethodCommand creates a Cobra command for a single service method.
//
//nolint:cyclop // unavoidable complexity for comprehensive method command building
func buildMethodCommand(serviceName string, _ reflect.Type, method reflect.Method, format *OutputFormat) *cobra.Command {
	methodName := method.Name
	kebabName := pascalToKebab(methodName)
	description := GetMethodDescription(serviceName, methodName)
	pattern := ClassifyMethod(method)

	// Determine if this method takes an option struct parameter.
	// Method signatures: receiver is index 0 (for reflect.Type.Method on pointer type).
	methodType := method.Type
	hasOpt := methodType.NumIn() == 2 //nolint:mnd // 2 = receiver + option param

	var (
		optType  reflect.Type
		optValue reflect.Value
	)

	if hasOpt {
		optType = methodType.In(1) // The option parameter type (should be *SomeOption)
		if optType.Kind() == reflect.Ptr {
			optType = optType.Elem()
		}
	}

	// Check if this is a streaming method.
	streamKey := serviceName + "." + methodName
	_, isStreaming := streamingMethods[streamKey]

	// Get response type for pagination support.
	var responseType reflect.Type

	if (pattern == PatternResponseBody || pattern == PatternSlice) && methodType.NumOut() == expectedTripleReturn {
		responseType = methodType.Out(0)
	}

	canPaginate := hasOpt && responseType != nil && hasPagination(optType) && responseHasPaging(responseType)

	cmd := &cobra.Command{ //nolint:exhaustruct // only setting fields relevant to method commands
		Use:   kebabName,
		Short: description,
		Long:  fmt.Sprintf("%s.%s — %s", serviceName, methodName, description),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runMethodCommand(cmd, serviceName, methodName, optValue, isStreaming, canPaginate, pattern, responseType, format)
		},
	}

	// Bind option struct fields as flags.
	if hasOpt {
		optValue = reflect.New(optType)
		BindFlags(cmd, optValue.Interface())
	}

	// Add --all flag for paginated methods.
	if canPaginate {
		cmd.Flags().Bool("all", false, "Fetch all pages of results (overrides --page and --page-size)")
	}

	return cmd
}

// runMethodCommand executes a service method command, handling streaming, pagination, and normal invocation.
func runMethodCommand(
	cmd *cobra.Command,
	serviceName, methodName string,
	optValue reflect.Value,
	isStreaming, canPaginate bool,
	pattern MethodReturnPattern,
	responseType reflect.Type,
	format *OutputFormat,
) error {
	client, err := clientFromContext(cmd)
	if err != nil {
		Logger().Error("failed to get SonarQube client", zap.Error(err))

		return err
	}

	// Get the service from the client.
	service := reflect.ValueOf(client).Elem().FieldByName(serviceName)
	if !service.IsValid() || service.IsNil() {
		err := fmt.Errorf("service %q not found on client", serviceName)
		Logger().Error("service not found", zap.String("service", serviceName))

		return err
	}

	if isStreaming {
		return InvokeStreamingMethod(service, methodName, optValue, os.Stdout)
	}

	// Check --all flag for pagination.
	allPages, _ := cmd.Flags().GetBool("all")

	if allPages && canPaginate {
		result, paginateErr := PaginateAll(service, methodName, optValue, pattern, responseType)
		if paginateErr != nil {
			Logger().Error("pagination failed",
				zap.String("service", serviceName),
				zap.String("method", methodName),
				zap.Error(paginateErr))

			return paginateErr
		}

		return FormatOutput(os.Stdout, result, *format)
	}

	hasOpt := optValue.IsValid()

	result, _, invokeErr := InvokeMethod(service, methodName, optValue, pattern, hasOpt) //nolint:bodyclose // CLI does not need to manage response body
	if invokeErr != nil {
		Logger().Error("method invocation failed",
			zap.String("service", serviceName),
			zap.String("method", methodName),
			zap.Error(invokeErr))

		return invokeErr
	}

	return FormatOutput(os.Stdout, result, *format)
}
