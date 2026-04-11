package main

import (
	"fmt"
	"os"

	"msgraphtool/internal/common/bootstrap"
	"msgraphtool/internal/common/logger"
	"msgraphtool/internal/common/version"
)

func main() {
	ctx, cancel := bootstrap.SetupSignalContext()
	defer cancel()

	config := parseAndConfigureFlags()

	if config.ShowVersion {
		fmt.Printf("pop3tool version %s\n", version.Get())
		fmt.Println("Part of gomailtesttool suite - https://github.com/ziembor/gomailtesttool")
		os.Exit(0)
	}

	if config.Action == "" {
		fmt.Fprintln(os.Stderr, "Error: -action is required")
		fmt.Fprintln(os.Stderr, "Use -help for usage information")
		os.Exit(1)
	}

	if err := validateConfiguration(config); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	slogLogger, csvLogger, err := bootstrap.InitLoggers("pop3tool", config.Action, config.VerboseMode, config.LogLevel, config.LogFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer csvLogger.Close()

	if err := executeAction(ctx, config, csvLogger, slogLogger); err != nil {
		logger.LogError(slogLogger, "Action failed", "action", config.Action, "error", err)
		os.Exit(1)
	}
}
