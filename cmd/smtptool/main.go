package main

import (
	"fmt"
	"os"

	"msgraphtool/internal/common/bootstrap"
	"msgraphtool/internal/common/logger"
	"msgraphtool/internal/common/version"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := bootstrap.SetupSignalContext()
	defer cancel()

	config := parseAndConfigureFlags()

	if config.ShowVersion {
		fmt.Printf("SMTP Connectivity Testing Tool - Version %s\n", version.Get())
		fmt.Println("Part of msgraphtool suite")
		fmt.Println("Repository: https://github.com/ziembor/gomailtesttool")
		return nil
	}

	if err := validateConfiguration(config); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	slogLogger, csvLogger, err := bootstrap.InitLoggers("smtptool", config.Action, config.VerboseMode, config.LogLevel, config.LogFormat)
	if err != nil {
		return err
	}
	defer csvLogger.Close()

	logger.LogInfo(slogLogger, "SMTP Connectivity Testing Tool started", "action", config.Action, "host", config.Host, "port", config.Port)

	if err := executeAction(ctx, config, csvLogger, slogLogger); err != nil {
		logger.LogError(slogLogger, "Action failed", "error", err)
		return err
	}

	logger.LogInfo(slogLogger, "Action completed successfully")
	return nil
}
