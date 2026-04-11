// Package main is a backward-compatibility shim for pop3tool.
// It translates the old -action=X flag style to the new Cobra subcommand style
// and delegates to the gomailtest binary.
//
// Old: pop3tool -action=testconnect -host=pop.example.com -port=110 ...
// New: gomailtest pop3 testconnect --host=pop.example.com --port=110 ...
//
// This shim will be removed in v3.1 once users have migrated to gomailtest.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	newArgs := translateArgs(os.Args[1:])

	gomailtest, err := findGomailtest()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: gomailtest binary not found in the same directory or PATH.")
		fmt.Fprintln(os.Stderr, "Install gomailtest alongside pop3tool, or add it to your PATH.")
		os.Exit(1)
	}

	cmd := exec.Command(gomailtest, newArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

// findGomailtest locates the gomailtest binary.
// It checks the same directory as this binary first, then falls back to PATH.
func findGomailtest() (string, error) {
	exePath, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(exePath), "gomailtest")
		if _, statErr := os.Stat(candidate); statErr == nil {
			return candidate, nil
		}
		// On Windows the binary has .exe extension
		candidateExe := candidate + ".exe"
		if _, statErr := os.Stat(candidateExe); statErr == nil {
			return candidateExe, nil
		}
	}

	return exec.LookPath("gomailtest")
}

// translateArgs converts old pop3tool flag-style arguments to the new Cobra
// subcommand style expected by gomailtest.
//
// Transformation:
//   - Extracts the action from -action=X or -action X (default: testconnect)
//   - Prepends ["pop3", actionName] to the remaining flags
//   - Converts single-dash long flags to double-dash (e.g. -host → --host)
//   - Drops -version flag (handled by gomailtest natively)
func translateArgs(args []string) []string {
	action := "testconnect" // default
	var rest []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Extract action
		if strings.HasPrefix(arg, "-action=") {
			action = strings.TrimPrefix(arg, "-action=")
			continue
		}
		if arg == "-action" && i+1 < len(args) {
			action = args[i+1]
			i++
			continue
		}

		// Drop -version flag (gomailtest has --version natively)
		if arg == "-version" || arg == "--version" {
			rest = append(rest, "--version")
			continue
		}

		// Convert single-dash long flags to double-dash.
		// Single-char short flags (-v) stay as-is; multi-char long flags get double-dash.
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			// Split on = to handle -flag=value style
			eqIdx := strings.Index(arg, "=")
			flagPart := arg
			valuePart := ""
			if eqIdx > 0 {
				flagPart = arg[:eqIdx]
				valuePart = arg[eqIdx:] // includes the "="
			}

			// Multi-character flag name gets double-dash
			flagName := strings.TrimPrefix(flagPart, "-")
			if len(flagName) > 1 {
				arg = "--" + flagName + valuePart
			}
		}

		rest = append(rest, arg)
	}
	return append([]string{"pop3", action}, rest...)
}
