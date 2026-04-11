// Package bootstrap provides shared startup helpers used by every cmd/*tool
// main.go: signal handling and logger initialization.
package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"msgraphtool/internal/common/logger"
)

// SetupSignalContext returns a context cancelled on SIGINT/SIGTERM.
// On signal it prints a shutdown notice to stdout and cancels the context.
func SetupSignalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\nReceived interrupt signal. Shutting down gracefully...")
		cancel()
	}()

	return ctx, cancel
}

// InitLoggers builds the structured slog logger and the file logger (CSV/JSON)
// from raw config primitives. Callers decide whether a non-nil error is fatal.
//
// The slog logger is always returned (never nil) so callers can log the
// failure. The file logger is nil when an error is returned.
func InitLoggers(toolName, action string, verbose bool, logLevel, logFormat string) (*slog.Logger, logger.Logger, error) {
	slogLogger := logger.SetupLogger(verbose, logLevel)

	format, err := logger.ParseLogFormat(logFormat)
	if err != nil {
		return slogLogger, nil, fmt.Errorf("invalid log format: %w", err)
	}

	fileLogger, err := logger.NewLogger(format, toolName, action)
	if err != nil {
		return slogLogger, nil, fmt.Errorf("failed to initialize file logger: %w", err)
	}

	return slogLogger, fileLogger, nil
}
