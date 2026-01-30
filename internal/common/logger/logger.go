package logger

import "fmt"

// Logger defines the common interface for all log formats (CSV, JSON, etc.).
// This allows tools to switch between formats without changing their logging code.
type Logger interface {
	// WriteHeader sets the column names for the logger.
	// For CSV, this writes a header row. For JSON, this stores column names.
	WriteHeader(columns []string) error

	// WriteRow writes a data row to the log.
	// The row values must match the columns set by WriteHeader.
	WriteRow(row []string) error

	// Close closes the logger and flushes any buffered data.
	Close() error

	// ShouldWriteHeader returns true if the log file is new and needs a header.
	ShouldWriteHeader() (bool, error)
}

// LogFormat represents the supported log formats.
type LogFormat string

const (
	// LogFormatCSV outputs CSV (Comma-Separated Values) format.
	LogFormatCSV LogFormat = "csv"

	// LogFormatJSON outputs JSONL (JSON Lines) format - one JSON object per line.
	LogFormatJSON LogFormat = "json"
)

// ParseLogFormat parses a log format string and returns the corresponding LogFormat.
// Returns an error if the format is not recognized.
func ParseLogFormat(format string) (LogFormat, error) {
	switch format {
	case "csv", "CSV":
		return LogFormatCSV, nil
	case "json", "JSON", "jsonl", "JSONL":
		return LogFormatJSON, nil
	default:
		return "", fmt.Errorf("unsupported log format: %s (valid options: csv, json)", format)
	}
}

// NewLogger creates a new logger based on the specified format.
// This is a factory function that returns the appropriate logger implementation.
//
// Parameters:
//   - format: The desired log format (csv, json)
//   - toolName: The tool name (e.g., "msgraphtool", "smtptool")
//   - action: The action being logged (e.g., "sendmail", "testauth")
//
// Returns a Logger interface that can be used for structured logging.
func NewLogger(format LogFormat, toolName, action string) (Logger, error) {
	switch format {
	case LogFormatCSV:
		return NewCSVLogger(toolName, action)
	case LogFormatJSON:
		return NewJSONLogger(toolName, action)
	default:
		return nil, fmt.Errorf("unsupported log format: %s", format)
	}
}
