package cli

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger //nolint:gochecknoglobals // global logger is a standard pattern

// initLogger creates a development-friendly logger suitable for CLI output.
func initLogger() {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// Write to stderr instead of stdout
	config.OutputPaths = []string{"stderr"}
	config.ErrorOutputPaths = []string{"stderr"}

	// Make output more CLI-friendly
	config.EncoderConfig.TimeKey = ""                                   // Don't include timestamps
	config.EncoderConfig.CallerKey = ""                                 // Don't include caller info
	config.EncoderConfig.StacktraceKey = ""                             // Don't include stack traces in console
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Colorize level (ERROR, WARN, etc.)
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder // Format durations as strings
	config.Encoding = "console"

	logger, err := config.Build()
	if err != nil {
		// Fallback to stderr if logger creation fails
		panic(err)
	}

	globalLogger = logger
}

// Logger returns the global logger instance.
func Logger() *zap.Logger {
	if globalLogger == nil {
		initLogger()
	}

	return globalLogger
}

// Sync flushes any buffered log entries to disk.
func Sync() error {
	if globalLogger != nil {
		err := globalLogger.Sync()
		if err != nil {
			return fmt.Errorf("failed to sync logger: %w", err)
		}

		return nil
	}

	return nil
}
