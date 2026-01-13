//go:build !integration
// +build !integration

package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"msgraphgolangtestingtool/internal/smtp/protocol"
)

// TestDebugLogCommand tests debug logging of SMTP commands
func TestDebugLogCommand(t *testing.T) {
	tests := []struct {
		name        string
		verboseMode bool
		command     string
		wantOutput  bool
		wantContain string
	}{
		{
			name:        "Verbose mode enabled - EHLO command",
			verboseMode: true,
			command:     "EHLO smtptool.local\r\n",
			wantOutput:  true,
			wantContain: ">>> EHLO smtptool.local",
		},
		{
			name:        "Verbose mode enabled - STARTTLS command",
			verboseMode: true,
			command:     "STARTTLS\r\n",
			wantOutput:  true,
			wantContain: ">>> STARTTLS",
		},
		{
			name:        "Verbose mode enabled - QUIT command",
			verboseMode: true,
			command:     "QUIT\r\n",
			wantOutput:  true,
			wantContain: ">>> QUIT",
		},
		{
			name:        "Verbose mode disabled - no output",
			verboseMode: false,
			command:     "EHLO smtptool.local\r\n",
			wantOutput:  false,
			wantContain: "",
		},
		{
			name:        "Verbose mode enabled - empty command",
			verboseMode: true,
			command:     "",
			wantOutput:  true,
			wantContain: ">>>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create SMTPClient with test config
			config := &Config{
				VerboseMode: tt.verboseMode,
			}
			client := &SMTPClient{
				config: config,
			}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call debug log command
			client.debugLogCommand(tt.command)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify output
			if tt.wantOutput {
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("debugLogCommand() output missing expected content\ngot:  %q\nwant: %q", output, tt.wantContain)
				}
				if !strings.Contains(output, ">>>") {
					t.Error("debugLogCommand() output missing >>> prefix")
				}
			} else {
				if output != "" {
					t.Errorf("debugLogCommand() produced output when verbose mode disabled: %q", output)
				}
			}
		})
	}
}

// TestDebugLogResponse tests debug logging of SMTP responses
func TestDebugLogResponse(t *testing.T) {
	tests := []struct {
		name        string
		verboseMode bool
		response    *protocol.SMTPResponse
		wantOutput  bool
		wantContain []string
	}{
		{
			name:        "Verbose mode enabled - single line response",
			verboseMode: true,
			response: &protocol.SMTPResponse{
				Code:    250,
				Message: "OK",
				Lines:   []string{"OK"},
			},
			wantOutput:  true,
			wantContain: []string{"<<< 250 OK"},
		},
		{
			name:        "Verbose mode enabled - multiline response",
			verboseMode: true,
			response: &protocol.SMTPResponse{
				Code:    250,
				Message: "smtp.example.com\nSTARTTLS\nAUTH PLAIN LOGIN",
				Lines:   []string{"smtp.example.com", "STARTTLS", "AUTH PLAIN LOGIN"},
			},
			wantOutput:  true,
			wantContain: []string{"<<< 250-smtp.example.com", "<<< 250-STARTTLS", "<<< 250 AUTH PLAIN LOGIN"},
		},
		{
			name:        "Verbose mode enabled - banner response",
			verboseMode: true,
			response: &protocol.SMTPResponse{
				Code:    220,
				Message: "smtp.gmail.com ESMTP",
				Lines:   []string{"smtp.gmail.com ESMTP"},
			},
			wantOutput:  true,
			wantContain: []string{"<<< 220 smtp.gmail.com ESMTP"},
		},
		{
			name:        "Verbose mode enabled - error response",
			verboseMode: true,
			response: &protocol.SMTPResponse{
				Code:    550,
				Message: "Mailbox not found",
				Lines:   []string{"Mailbox not found"},
			},
			wantOutput:  true,
			wantContain: []string{"<<< 550 Mailbox not found"},
		},
		{
			name:        "Verbose mode disabled - no output",
			verboseMode: false,
			response: &protocol.SMTPResponse{
				Code:    250,
				Message: "OK",
				Lines:   []string{"OK"},
			},
			wantOutput:  false,
			wantContain: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create SMTPClient with test config
			config := &Config{
				VerboseMode: tt.verboseMode,
			}
			client := &SMTPClient{
				config: config,
			}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call debug log response
			client.debugLogResponse(tt.response)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify output
			if tt.wantOutput {
				for _, expectedContent := range tt.wantContain {
					if !strings.Contains(output, expectedContent) {
						t.Errorf("debugLogResponse() output missing expected content\ngot:  %q\nwant: %q", output, expectedContent)
					}
				}
				if !strings.Contains(output, "<<<") {
					t.Error("debugLogResponse() output missing <<< prefix")
				}
			} else {
				if output != "" {
					t.Errorf("debugLogResponse() produced output when verbose mode disabled: %q", output)
				}
			}
		})
	}
}

// TestDebugLogMessage tests debug logging of informational messages
func TestDebugLogMessage(t *testing.T) {
	tests := []struct {
		name        string
		verboseMode bool
		message     string
		wantOutput  bool
		wantContain string
	}{
		{
			name:        "Verbose mode enabled - TLS handshake message",
			verboseMode: true,
			message:     "Performing TLS handshake...",
			wantOutput:  true,
			wantContain: "... Performing TLS handshake...",
		},
		{
			name:        "Verbose mode enabled - auth mechanism message",
			verboseMode: true,
			message:     "Starting authentication with mechanism: PLAIN",
			wantOutput:  true,
			wantContain: "... Starting authentication with mechanism: PLAIN",
		},
		{
			name:        "Verbose mode disabled - no output",
			verboseMode: false,
			message:     "This should not appear",
			wantOutput:  false,
			wantContain: "",
		},
		{
			name:        "Verbose mode enabled - empty message",
			verboseMode: true,
			message:     "",
			wantOutput:  true,
			wantContain: "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create SMTPClient with test config
			config := &Config{
				VerboseMode: tt.verboseMode,
			}
			client := &SMTPClient{
				config: config,
			}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call debug log message
			client.debugLogMessage(tt.message)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify output
			if tt.wantOutput {
				if !strings.Contains(output, tt.wantContain) {
					t.Errorf("debugLogMessage() output missing expected content\ngot:  %q\nwant: %q", output, tt.wantContain)
				}
				if !strings.Contains(output, "...") {
					t.Error("debugLogMessage() output missing ... prefix")
				}
			} else {
				if output != "" {
					t.Errorf("debugLogMessage() produced output when verbose mode disabled: %q", output)
				}
			}
		})
	}
}

// TestDebugLogging_NilSafety tests that debug logging handles nil values safely
func TestDebugLogging_NilSafety(t *testing.T) {
	t.Run("debugLogResponse with nil response", func(t *testing.T) {
		config := &Config{
			VerboseMode: true,
		}
		client := &SMTPClient{
			config: config,
		}

		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("debugLogResponse() panicked with nil response: %v", r)
			}
		}()

		client.debugLogResponse(nil)
	})

	t.Run("debugLogCommand with nil config", func(t *testing.T) {
		client := &SMTPClient{
			config: nil,
		}

		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("debugLogCommand() panicked with nil config: %v", r)
			}
		}()

		client.debugLogCommand("EHLO test\r\n")
	})
}

// TestDebugLogging_MultilineResponseFormatting tests correct formatting of multiline SMTP responses
func TestDebugLogging_MultilineResponseFormatting(t *testing.T) {
	config := &Config{
		VerboseMode: true,
	}
	client := &SMTPClient{
		config: config,
	}

	// Create a typical EHLO response with multiple capabilities
	response := &protocol.SMTPResponse{
		Code: 250,
		Message: "smtp.example.com\nSIZE 35882577\n8BITMIME\nSTARTTLS\nAUTH PLAIN LOGIN",
		Lines: []string{
			"smtp.example.com",
			"SIZE 35882577",
			"8BITMIME",
			"STARTTLS",
			"AUTH PLAIN LOGIN",
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	client.debugLogResponse(response)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify multiline format
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 5 {
		t.Errorf("Expected 5 output lines for multiline response, got %d", len(lines))
	}

	// First 4 lines should have hyphen after code
	for i := 0; i < 4; i++ {
		if !strings.HasPrefix(lines[i], "<<< 250-") {
			t.Errorf("Line %d should start with '<<< 250-', got: %s", i, lines[i])
		}
	}

	// Last line should have space after code
	if !strings.HasPrefix(lines[4], "<<< 250 ") {
		t.Errorf("Last line should start with '<<< 250 ', got: %s", lines[4])
	}
}
