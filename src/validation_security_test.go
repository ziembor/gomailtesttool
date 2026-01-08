//go:build !integration
// +build !integration

package main

import (
	"os"
	"strings"
	"testing"
)

// TestValidateMessageID_MaxLength tests the 998 character limit validation
func TestValidateMessageID_MaxLength(t *testing.T) {
	tests := []struct {
		name    string
		msgID   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "exactly 998 characters is valid",
			msgID:   "<" + strings.Repeat("a", 996) + ">", // 996 + 2 brackets = 998
			wantErr: false,
		},
		{
			name:    "999 characters exceeds limit",
			msgID:   "<" + strings.Repeat("a", 997) + ">", // 997 + 2 brackets = 999
			wantErr: true,
			errMsg:  "exceeds maximum length",
		},
		{
			name:    "1000 characters exceeds limit",
			msgID:   "<" + strings.Repeat("b", 998) + ">", // 998 + 2 brackets = 1000
			wantErr: true,
			errMsg:  "exceeds maximum length",
		},
		{
			name:    "very long message ID exceeds limit",
			msgID:   "<" + strings.Repeat("x", 2000) + ">",
			wantErr: true,
			errMsg:  "exceeds maximum length of 998 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessageID(tt.msgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMessageID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateMessageID() error = %q, should contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestValidateMessageID_AllODataOperators tests all OData operators are blocked
func TestValidateMessageID_AllODataOperators(t *testing.T) {
	// Test all operators to ensure 100% coverage of the loop
	operators := []string{" or ", " and ", " eq ", " ne ", " lt ", " gt ", " le ", " ge ", " not "}

	for _, op := range operators {
		t.Run("reject_"+strings.TrimSpace(op), func(t *testing.T) {
			// Create message ID with the operator embedded
			msgID := "<test" + op + "value@domain.com>"
			err := validateMessageID(msgID)

			if err == nil {
				t.Errorf("validateMessageID(%q) should reject OData operator %q", msgID, op)
				return
			}

			if !strings.Contains(err.Error(), "OData operators") {
				t.Errorf("validateMessageID(%q) error = %q, should mention OData operators", msgID, err.Error())
			}
		})
	}
}

// TestValidateMessageID_CaseSensitivity tests that operators are detected case-insensitively
func TestValidateMessageID_CaseSensitivity(t *testing.T) {
	tests := []struct {
		name  string
		msgID string
	}{
		{"uppercase OR", "<test OR value@domain.com>"},
		{"mixed case And", "<test And value@domain.com>"},
		{"uppercase EQ", "<test EQ value@domain.com>"},
		{"mixed case Ne", "<test Ne value@domain.com>"},
		{"uppercase LT", "<test LT value@domain.com>"},
		{"mixed case Gt", "<test Gt value@domain.com>"},
		{"uppercase LE", "<test LE value@domain.com>"},
		{"mixed case Ge", "<test Ge value@domain.com>"},
		{"uppercase NOT", "<test NOT value@domain.com>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessageID(tt.msgID)
			if err == nil {
				t.Errorf("validateMessageID(%q) should reject case-insensitive OData operator", tt.msgID)
			}
		})
	}
}

// TestValidateConfiguration_BodyTemplate tests body template file validation
func TestValidateConfiguration_BodyTemplate(t *testing.T) {
	// Create a temporary body template file
	tmpTemplate, err := os.CreateTemp("", "template-*.html")
	if err != nil {
		t.Fatalf("Failed to create temp template file: %v", err)
	}
	defer os.Remove(tmpTemplate.Name())
	tmpTemplate.WriteString("<html><body>{{.Content}}</body></html>")
	tmpTemplate.Close()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid body template file",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendMail,
				BodyTemplate: tmpTemplate.Name(),
				OutputFormat: "text",
			},
			wantErr: false,
		},
		{
			name: "body template file does not exist",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendMail,
				BodyTemplate: "/nonexistent/template.html",
				OutputFormat: "text",
			},
			wantErr: true,
			errMsg:  "Body template file",
		},
		{
			name: "body template with path traversal",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendMail,
				BodyTemplate: "../../etc/passwd",
				OutputFormat: "text",
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
					t.Errorf("validateConfiguration() error = %q, should contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestValidateConfiguration_EmailLists tests CC and BCC validation
func TestValidateConfiguration_EmailLists(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid CC recipients",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendMail,
				Cc:           stringSlice{"cc1@example.com", "cc2@example.com"},
				OutputFormat: "text",
			},
			wantErr: false,
		},
		{
			name: "invalid CC recipient",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendMail,
				Cc:           stringSlice{"valid@example.com", "invalid-email"},
				OutputFormat: "text",
			},
			wantErr: true,
			errMsg:  "CC recipients",
		},
		{
			name: "valid BCC recipients",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendMail,
				Bcc:          stringSlice{"bcc1@example.com", "bcc2@example.com"},
				OutputFormat: "text",
			},
			wantErr: false,
		},
		{
			name: "invalid BCC recipient",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendMail,
				Bcc:          stringSlice{"bcc@example.com", "not-an-email"},
				OutputFormat: "text",
			},
			wantErr: true,
			errMsg:  "BCC recipients",
		},
		{
			name: "valid combination of To, CC, and BCC",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendMail,
				To:           stringSlice{"to@example.com"},
				Cc:           stringSlice{"cc@example.com"},
				Bcc:          stringSlice{"bcc@example.com"},
				OutputFormat: "text",
			},
			wantErr: false,
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
					t.Errorf("validateConfiguration() error = %q, should contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestValidateFilePath_PermissionDenied tests permission error handling
// Note: This test creates a file and attempts to make it unreadable, which may not work on all systems
func TestValidateFilePath_ErrorHandling(t *testing.T) {
	// Test with a path that would cause filepath.Abs to succeed but has other issues
	// We can't easily trigger permission denied without platform-specific code,
	// so we'll test the error path coverage by testing edge cases

	tests := []struct {
		name      string
		path      string
		fieldName string
		setup     func() (string, func()) // Returns path and cleanup function
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty path is allowed",
			path:      "",
			fieldName: "Test file",
			setup:     nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cleanup func()
			testPath := tt.path

			if tt.setup != nil {
				testPath, cleanup = tt.setup()
				if cleanup != nil {
					defer cleanup()
				}
			}

			err := validateFilePath(testPath, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateFilePath() error = %q, should contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestValidateConfiguration_Comprehensive tests comprehensive validation coverage
func TestValidateConfiguration_Comprehensive(t *testing.T) {
	// Create temporary files for testing
	tmpAttach, err := os.CreateTemp("", "attach-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp attachment: %v", err)
	}
	defer os.Remove(tmpAttach.Name())
	tmpAttach.Close()

	tmpTemplate, err := os.CreateTemp("", "template-*.html")
	if err != nil {
		t.Fatalf("Failed to create temp template: %v", err)
	}
	defer os.Remove(tmpTemplate.Name())
	tmpTemplate.Close()

	tmpPfx, err := os.CreateTemp("", "cert-*.pfx")
	if err != nil {
		t.Fatalf("Failed to create temp pfx: %v", err)
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
			name: "comprehensive valid configuration with all fields",
			config: &Config{
				TenantID:        "12345678-1234-1234-1234-123456789012",
				ClientID:        "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:         "sender@example.com",
				Secret:          "test-secret",
				Action:          ActionSendMail,
				To:              stringSlice{"to1@example.com", "to2@example.com"},
				Cc:              stringSlice{"cc1@example.com", "cc2@example.com"},
				Bcc:             stringSlice{"bcc1@example.com"},
				AttachmentFiles: stringSlice{tmpAttach.Name()},
				BodyTemplate:    tmpTemplate.Name(),
				StartTime:       "2026-01-15T10:00:00Z",
				EndTime:         "2026-01-15T11:00:00Z",
				OutputFormat:    "json",
			},
			wantErr: false,
		},
		{
			name: "PFX certificate authentication",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				PfxPath:      tmpPfx.Name(),
				PfxPass:      "password",
				Action:       ActionGetInbox,
				OutputFormat: "text",
			},
			wantErr: false,
		},
		{
			name: "invalid PFX path",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				PfxPath:      "/nonexistent/cert.pfx",
				PfxPass:      "password",
				Action:       ActionGetInbox,
				OutputFormat: "text",
			},
			wantErr: true,
			errMsg:  "PFX certificate file",
		},
		{
			name: "start time validation error",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendInvite,
				StartTime:    "invalid-time",
				OutputFormat: "text",
			},
			wantErr: true,
			errMsg:  "Start time",
		},
		{
			name: "end time validation error",
			config: &Config{
				TenantID:     "12345678-1234-1234-1234-123456789012",
				ClientID:     "abcdefab-1234-1234-1234-abcdefabcdef",
				Mailbox:      "test@example.com",
				Secret:       "test-secret",
				Action:       ActionSendInvite,
				EndTime:      "not-a-valid-time",
				OutputFormat: "text",
			},
			wantErr: true,
			errMsg:  "End time",
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
					t.Errorf("validateConfiguration() error = %q, should contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestValidateFilePath_RelativePathHandling tests relative path handling edge cases
func TestValidateFilePath_RelativePathHandling(t *testing.T) {
	// Create a temp file for absolute path testing
	tmpFile, err := os.CreateTemp("", "reltest-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name      string
		path      string
		fieldName string
		wantErr   bool
	}{
		{
			name:      "absolute path to file succeeds",
			path:      tmpFile.Name(),
			fieldName: "Test file",
			wantErr:   false,
		},
		{
			name:      "non-existent relative path fails",
			path:      "nonexistent-file-12345.txt",
			fieldName: "Test file",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.path, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFilePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}
