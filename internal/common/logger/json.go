package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// JSONLogger handles JSON logging operations with periodic buffering.
// Outputs JSONL (JSON Lines) format - one JSON object per line.
type JSONLogger struct {
	file       *os.File
	toolName   string    // Tool name for filename
	action     string    // Action being performed
	columns    []string  // Column names for structured output
	rowCount   int       // Number of rows written since last flush
	lastFlush  time.Time // Time of last flush
	flushEvery int       // Flush every N rows
}

// NewJSONLogger creates a new JSON logger for the specified tool and action.
// The toolName parameter differentiates between tools (e.g., "msgraphgolangtestingtool" or "smtptool").
// Filename pattern: %TEMP%\_{toolName}_{action}_{date}.jsonl
//
// Examples:
//   - _msgraphgolangtestingtool_sendmail_2026-01-09.jsonl
//   - _smtptool_teststarttls_2026-01-09.jsonl
//
// The output format is JSONL (JSON Lines): one JSON object per line.
// Each line is a valid JSON object with timestamp and structured fields.
func NewJSONLogger(toolName, action string) (*JSONLogger, error) {
	// Get temp directory
	tempDir := os.TempDir()

	// Create filename with tool name, action, and current date
	dateStr := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("_%s_%s_%s.jsonl", toolName, action, dateStr)
	filePath := filepath.Join(tempDir, fileName)

	// Open or create file (append mode) with restrictive permissions (0600)
	// This ensures only the owner can read/write the file, protecting sensitive data
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("could not create JSON log file: %w", err)
	}

	// Apply platform-specific restrictive permissions for additional security
	if err := setRestrictivePermissions(file, filePath); err != nil {
		// Log warning but continue - file creation succeeded
		log.Printf("Warning: Failed to set restrictive permissions on log file %s: %v", filePath, err)
	}

	logger := &JSONLogger{
		file:       file,
		toolName:   toolName,
		action:     action,
		columns:    nil, // Set by WriteHeader
		rowCount:   0,
		lastFlush:  time.Now(),
		flushEvery: 10, // Flush every 10 rows or on close
	}

	fmt.Printf("Logging to: %s\n\n", filePath)
	return logger, nil
}

// WriteHeader sets the column names for structured JSON output.
// Unlike CSV, JSONL doesn't need a header row, but we store column names
// to create properly structured JSON objects.
func (l *JSONLogger) WriteHeader(columns []string) error {
	// Store column names for use in WriteRow
	l.columns = columns
	return nil
}

// WriteRow writes a JSON object to the log file.
// Each row becomes a JSON object with timestamp and fields named according to the header.
// The output format is JSONL: one JSON object per line.
//
// Example output:
// {"timestamp":"2026-01-09 14:30:45","action":"sendmail","status":"SUCCESS","server":"smtp.example.com"}
func (l *JSONLogger) WriteRow(row []string) error {
	if l.file == nil {
		return fmt.Errorf("JSON file is not initialized")
	}

	if l.columns == nil {
		return fmt.Errorf("WriteHeader must be called before WriteRow")
	}

	if len(row) != len(l.columns) {
		return fmt.Errorf("row length (%d) does not match header length (%d)", len(row), len(l.columns))
	}

	// Create JSON object with timestamp and fields
	obj := make(map[string]string)
	obj["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	// Map column names to values
	for i, colName := range l.columns {
		obj[colName] = row[i]
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write JSON line
	if _, err := l.file.Write(jsonBytes); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	// Write newline
	if _, err := l.file.Write([]byte("\n")); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	l.rowCount++

	// Flush every N rows or every 5 seconds
	if l.rowCount%l.flushEvery == 0 || time.Since(l.lastFlush) > 5*time.Second {
		if err := l.file.Sync(); err != nil {
			return fmt.Errorf("failed to flush JSON: %w", err)
		}
		l.lastFlush = time.Now()
	}

	return nil
}

// Close closes the JSON file, ensuring all buffered data is flushed.
// Always call this method when done logging to prevent data loss.
func (l *JSONLogger) Close() error {
	if l.file != nil {
		if err := l.file.Sync(); err != nil {
			return fmt.Errorf("error flushing JSON on close: %w", err)
		}
		return l.file.Close()
	}
	return nil
}

// ShouldWriteHeader checks if the JSON file is new (empty) and needs a header.
// For JSON logs, this always returns true since we need to set column names.
// However, JSONL format doesn't actually write a header row.
func (l *JSONLogger) ShouldWriteHeader() (bool, error) {
	fileInfo, err := l.file.Stat()
	if err != nil {
		return false, fmt.Errorf("could not stat JSON file: %w", err)
	}
	return fileInfo.Size() == 0, nil
}
