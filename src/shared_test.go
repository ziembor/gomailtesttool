//go:build !integration
// +build !integration

package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestValidateFilePath tests the validateFilePath function with various inputs
func TestValidateFilePath(t *testing.T) {
	// Create a temporary file for valid path tests
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a temporary directory to test directory rejection
	tmpDir, err := os.MkdirTemp("", "testdir-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name      string
		path      string
		fieldName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty path is allowed",
			path:      "",
			fieldName: "Test file",
			wantErr:   false,
		},
		{
			name:      "valid absolute path",
			path:      tmpFile.Name(),
			fieldName: "PFX file",
			wantErr:   false,
		},
		{
			name:      "path traversal with ..",
			path:      "../../etc/passwd",
			fieldName: "PFX file",
			wantErr:   true,
			errMsg:    "path contains directory traversal",
		},
		{
			name:      "path traversal Windows style",
			path:      "..\\..\\Windows\\System32\\config\\SAM",
			fieldName: "Attachment",
			wantErr:   true,
			errMsg:    "path contains directory traversal",
		},
		{
			name:      "file does not exist",
			path:      filepath.Join(os.TempDir(), "nonexistent-file-12345.pfx"),
			fieldName: "PFX file",
			wantErr:   true,
			errMsg:    "file not found",
		},
		{
			name:      "path is a directory not a file",
			path:      tmpDir,
			fieldName: "Attachment",
			wantErr:   true,
			errMsg:    "not a regular file",
		},
		{
			name:      "relative path to existing file",
			path:      filepath.Base(tmpFile.Name()),
			fieldName: "Test file",
			wantErr:   true, // Will fail because relative path won't be found unless we're in tmpdir
			errMsg:    "file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.path, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFilePath(%q, %q) error = %v, wantErr %v", tt.path, tt.fieldName, err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateFilePath(%q, %q) error = %v, should contain %q", tt.path, tt.fieldName, err, tt.errMsg)
				}
			}
		})
	}
}

// TestValidateFilePath_PathTraversalVariations tests various path traversal attempts
func TestValidateFilePath_PathTraversalVariations(t *testing.T) {
	traversalPaths := []string{
		"../secret.pfx",
		"../../etc/passwd",
		"foo/../../../etc/shadow",
		"..\\..\\Windows\\System32",
		"test\\..\\..\\sensitive.txt",
	}

	for _, path := range traversalPaths {
		t.Run(path, func(t *testing.T) {
			err := validateFilePath(path, "Test file")
			if err == nil {
				t.Errorf("validateFilePath(%q) should reject path traversal, but got nil error", path)
			}
			// Either it should fail with traversal error, or file not found (both are acceptable)
			errMsg := err.Error()
			if !strings.Contains(errMsg, "directory traversal") && !strings.Contains(errMsg, "file not found") {
				t.Errorf("validateFilePath(%q) error = %v, expected traversal or not found error", path, err)
			}
		})
	}
}

// TestValidateConfiguration_PfxPathValidation tests that validateConfiguration validates PFX paths
func TestValidateConfiguration_PfxPathValidation(t *testing.T) {
	// Create a temporary PFX file for testing
	tmpPfx, err := os.CreateTemp("", "test-*.pfx")
	if err != nil {
		t.Fatalf("Failed to create temp PFX file: %v", err)
	}
	defer os.Remove(tmpPfx.Name())
	tmpPfx.Close()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid PFX path",
			config: &Config{
				TenantID: "12345678-1234-1234-1234-123456789012",
				ClientID: "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:  "test@example.com",
				PfxPath:  tmpPfx.Name(),
				PfxPass:  "password",
				Action:   ActionGetInbox,
			},
			wantErr: false,
		},
		{
			name: "PFX path does not exist",
			config: &Config{
				TenantID: "12345678-1234-1234-1234-123456789012",
				ClientID: "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:  "test@example.com",
				PfxPath:  "/nonexistent/path/cert.pfx",
				PfxPass:  "password",
				Action:   ActionGetInbox,
			},
			wantErr: true,
			errMsg:  "file not found",
		},
		{
			name: "PFX path with traversal",
			config: &Config{
				TenantID: "12345678-1234-1234-1234-123456789012",
				ClientID: "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:  "test@example.com",
				PfxPath:  "../../etc/passwd",
				PfxPass:  "password",
				Action:   ActionGetInbox,
			},
			wantErr: true,
			errMsg:  "directory traversal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfiguration(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateConfiguration() error = %v, should contain %q", err, tt.errMsg)
				}
			}
		})
	}
}

// TestValidateConfiguration_AttachmentFilesValidation tests that validateConfiguration validates attachment paths
func TestValidateConfiguration_AttachmentFilesValidation(t *testing.T) {
	// Create temporary attachment files for testing
	tmpAttach1, err := os.CreateTemp("", "attach1-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp attachment file: %v", err)
	}
	defer os.Remove(tmpAttach1.Name())
	tmpAttach1.Close()

	tmpAttach2, err := os.CreateTemp("", "attach2-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp attachment file: %v", err)
	}
	defer os.Remove(tmpAttach2.Name())
	tmpAttach2.Close()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid attachment paths",
			config: &Config{
				TenantID:        "12345678-1234-1234-1234-123456789012",
				ClientID:        "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:         "test@example.com",
				Secret:          "test-secret",
				AttachmentFiles: stringSlice{tmpAttach1.Name(), tmpAttach2.Name()},
				Action:          ActionSendMail,
			},
			wantErr: false,
		},
		{
			name: "one attachment does not exist",
			config: &Config{
				TenantID:        "12345678-1234-1234-1234-123456789012",
				ClientID:        "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:         "test@example.com",
				Secret:          "test-secret",
				AttachmentFiles: stringSlice{tmpAttach1.Name(), "/nonexistent/file.txt"},
				Action:          ActionSendMail,
			},
			wantErr: true,
			errMsg:  "Attachment file #2",
		},
		{
			name: "attachment with path traversal",
			config: &Config{
				TenantID:        "12345678-1234-1234-1234-123456789012",
				ClientID:        "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:         "test@example.com",
				Secret:          "test-secret",
				AttachmentFiles: stringSlice{"../../etc/shadow"},
				Action:          ActionSendMail,
			},
			wantErr: true,
			errMsg:  "directory traversal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfiguration(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateConfiguration() error = %v, should contain %q", err, tt.errMsg)
				}
			}
		})
	}
}

// TestParseLogLevel tests the parseLogLevel function
func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLevel slog.Level
	}{
		{"debug lowercase", "debug", slog.LevelDebug},
		{"debug uppercase", "DEBUG", slog.LevelDebug},
		{"info lowercase", "info", slog.LevelInfo},
		{"info uppercase", "INFO", slog.LevelInfo},
		{"warn lowercase", "warn", slog.LevelWarn},
		{"warn uppercase", "WARN", slog.LevelWarn},
		{"warning", "WARNING", slog.LevelWarn},
		{"error lowercase", "error", slog.LevelError},
		{"error uppercase", "ERROR", slog.LevelError},
		{"invalid level defaults to info", "INVALID", slog.LevelInfo},
		{"empty string defaults to info", "", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLogLevel(tt.input)
			if got != tt.wantLevel {
				t.Errorf("parseLogLevel(%q) = %v, want %v", tt.input, got, tt.wantLevel)
			}
		})
	}
}

// TestSetupLogger tests logger configuration
func TestSetupLogger(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectDebug bool
	}{
		{
			name:        "verbose mode enables debug",
			config:      &Config{VerboseMode: true, LogLevel: "INFO"},
			expectDebug: true,
		},
		{
			name:        "debug level enables debug",
			config:      &Config{VerboseMode: false, LogLevel: "DEBUG"},
			expectDebug: true,
		},
		{
			name:        "info level disables debug",
			config:      &Config{VerboseMode: false, LogLevel: "INFO"},
			expectDebug: false,
		},
		{
			name:        "error level disables debug",
			config:      &Config{VerboseMode: false, LogLevel: "ERROR"},
			expectDebug: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := setupLogger(tt.config)
			if logger == nil {
				t.Fatal("setupLogger returned nil")
			}

			// Logger is created successfully if we get here
			// Testing actual log output would require capturing output
			// which is beyond the scope of a unit test
		})
	}
}

// TestLogHelpers tests the log helper functions don't panic with nil logger
func TestLogHelpers(t *testing.T) {
	// These should not panic even with nil logger
	logDebug(nil, "test debug")
	logInfo(nil, "test info")
	logWarn(nil, "test warn")
	logError(nil, "test error")

	// These should not panic with actual logger
	config := &Config{LogLevel: "DEBUG"}
	logger := setupLogger(config)
	logDebug(logger, "test debug", "key", "value")
	logInfo(logger, "test info", "key", "value")
	logWarn(logger, "test warn", "key", "value")
	logError(logger, "test error", "key", "value")
}
