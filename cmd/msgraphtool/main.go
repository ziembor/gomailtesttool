//go:build !integration
// +build !integration

// Package main provides a portable CLI tool for interacting with Microsoft Graph API
// to manage Exchange Online (EXO) mailboxes. The tool supports sending emails,
// creating calendar events, and retrieving inbox messages and calendar events.
//
// Authentication methods supported:
//   - Client Secret: Standard App Registration secret
//   - PFX Certificate: Certificate file with private key
//   - Windows Certificate Store: Thumbprint-based certificate retrieval (Windows only)
//
// All operations are automatically logged to action-specific CSV files in the
// system temp directory for audit and troubleshooting purposes.
//
// Example usage:
//
//	msgraphtool.exe -tenantid "..." -clientid "..." -secret "..." -mailbox "user@example.com" -action sendmail
//
// Version information is embedded from the VERSION file at compile time using go:embed.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"msgraphtool/internal/common/bootstrap"
	"msgraphtool/internal/common/logger"
	"msgraphtool/internal/common/version"
)

func main() {
	// Handle -completion flag FIRST, before anything else runs
	// This ensures only completion script is output, all other flags are ignored
	for i, arg := range os.Args {
		if arg == "-completion" && i+1 < len(os.Args) {
			shellType := os.Args[i+1]
			if shellType == "bash" {
				fmt.Print(generateBashCompletion())
				os.Exit(0)
			} else if shellType == "powershell" {
				fmt.Print(generatePowerShellCompletion())
				os.Exit(0)
			} else {
				fmt.Fprintf(os.Stderr, "Error: Invalid completion shell type '%s'\n", shellType)
				fmt.Fprintf(os.Stderr, "Valid options: bash, powershell\n\n")
				fmt.Fprintf(os.Stderr, "Usage:\n")
				fmt.Fprintf(os.Stderr, "  %s -completion bash > msgraphtool-completion.bash\n", os.Args[0])
				fmt.Fprintf(os.Stderr, "  %s -completion powershell > msgraphtool-completion.ps1\n", os.Args[0])
				os.Exit(1)
			}
		}
	}

	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// initializeServices sets up file logging and proxy configuration based on
// the provided configuration. If file logger initialization fails, a warning
// is logged but execution continues with a nil logger.
func initializeServices(config *Config) (logger.Logger, error) {
	_, csvLogger, err := bootstrap.InitLoggers("msgraphtool", config.Action, config.VerboseMode, config.LogLevel, config.LogFormat)
	if err != nil {
		log.Printf("Warning: Could not initialize file logging: %v", err)
		csvLogger = nil // Continue without logging
	}

	if config.ProxyURL != "" {
		os.Setenv("HTTP_PROXY", config.ProxyURL)
		os.Setenv("HTTPS_PROXY", config.ProxyURL)
		fmt.Printf("Using proxy: %s\n", config.ProxyURL)
	}

	return csvLogger, nil
}

// run is the main application entry point that orchestrates the tool's execution flow.
// It performs the following steps:
//  1. Sets up graceful shutdown handling for interrupt signals
//  2. Parses and validates configuration from flags and environment variables
//  3. Initializes services (CSV logging, proxy configuration)
//  4. Creates Microsoft Graph SDK client with appropriate authentication
//  5. Executes the requested action (getevents, sendmail, sendinvite, getinbox)
//
// Returns an error if any step fails, nil on successful completion.
func run() error {
	// 1. Setup signal handling for graceful shutdown
	ctx, cancel := bootstrap.SetupSignalContext()
	defer cancel()

	// 2. Parse command-line flags and apply environment variables
	config := parseAndConfigureFlags()

	// 3. Handle version flag early exit
	if config.ShowVersion {
		fmt.Printf("Microsoft Graph EXO Mails/Calendar Golang Testing Tool - Version %s\n", version.Get())
		return nil
	}

	// 4. Validate configuration
	if err := validateConfiguration(config); err != nil {
		flag.Usage()
		return fmt.Errorf("validation failed: %w", err)
	}

	// 5. Setup structured logger
	slogger := logger.SetupLogger(config.VerboseMode, config.LogLevel)
	logger.LogInfo(slogger, "Application starting", "version", version.Get(), "action", config.Action)

	// Load body template if provided (validation already done in step 4)
	if config.BodyTemplate != "" {
		content, err := os.ReadFile(config.BodyTemplate)
		if err != nil {
			return fmt.Errorf("failed to read body template file: %w", err)
		}
		config.BodyHTML = string(content)
		slogger.Info("Loaded email body from template", "path", config.BodyTemplate, "size", len(config.BodyHTML))
	}

	// 6. Initialize services (CSV logging and proxy)
	csvLogger, _ := initializeServices(config)
	if csvLogger != nil {
		defer csvLogger.Close()
	}

	// 7. Setup Microsoft Graph client
	client, err := setupGraphClient(ctx, config, slogger)
	if err != nil {
		return err
	}

	// 8. Execute the requested action
	return executeAction(ctx, client, config, csvLogger)
}
