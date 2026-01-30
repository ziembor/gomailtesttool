//go:build !integration
// +build !integration

package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestNewCSVLogger tests creating CSV loggers for different actions
func TestNewCSVLogger(t *testing.T) {
	tests := []struct {
		name           string
		action         string
		wantErr        bool
		expectedHeader []string
	}{
		{
			name:           "getevents action",
			action:         ActionGetEvents,
			wantErr:        false,
			expectedHeader: []string{"Timestamp", "Action", "Status", "Mailbox", "Event Subject", "Event ID"},
		},
		{
			name:           "sendmail action",
			action:         ActionSendMail,
			wantErr:        false,
			expectedHeader: []string{"Timestamp", "Action", "Status", "Mailbox", "To", "CC", "BCC", "Subject", "Body Type", "Attachments"},
		},
		{
			name:           "sendinvite action",
			action:         ActionSendInvite,
			wantErr:        false,
			expectedHeader: []string{"Timestamp", "Action", "Status", "Mailbox", "Subject", "Start Time", "End Time", "Event ID"},
		},
		{
			name:           "getinbox action",
			action:         ActionGetInbox,
			wantErr:        false,
			expectedHeader: []string{"Timestamp", "Action", "Status", "Mailbox", "Subject", "From", "To", "Received DateTime"},
		},
		{
			name:           "getschedule action",
			action:         ActionGetSchedule,
			wantErr:        false,
			expectedHeader: []string{"Timestamp", "Action", "Status", "Mailbox", "Recipient", "Check DateTime", "Availability View"},
		},
		{
			name:           "exportinbox action",
			action:         ActionExportInbox,
			wantErr:        false,
			expectedHeader: []string{"Timestamp", "Action", "Status", "Details"},
		},
		{
			name:           "searchandexport action",
			action:         ActionSearchAndExport,
			wantErr:        false,
			expectedHeader: []string{"Timestamp", "Action", "Status", "Details"},
		},
		{
			name:           "unknown action defaults to generic header",
			action:         "unknownaction",
			wantErr:        false,
			expectedHeader: []string{"Timestamp", "Action", "Status", "Details"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger
			logger, err := NewCSVLogger(tt.action)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewCSVLogger(%q) error = %v, wantErr %v", tt.action, err, tt.wantErr)
				return
			}

			if err != nil {
				return // Expected error, test passed
			}

			// Verify logger is not nil
			if logger == nil {
				t.Fatal("NewCSVLogger() returned nil logger with no error")
			}

			// Close logger to flush data
			if err := logger.Close(); err != nil {
				t.Errorf("Failed to close logger: %v", err)
			}

			// Verify file was created with expected name pattern
			tempDir := os.TempDir()
			dateStr := time.Now().Format("2006-01-02")
			expectedFileName := filepath.Join(tempDir, "_msgraphtool_"+tt.action+"_"+dateStr+".csv")

			// Check file exists
			fileInfo, err := os.Stat(expectedFileName)
			if err != nil {
				t.Errorf("CSV file not created at expected path: %s, error: %v", expectedFileName, err)
				return
			}

			// Verify file is not empty (should have at least header)
			if fileInfo.Size() == 0 {
				t.Error("CSV file is empty, expected at least header row")
			}

			// Read and verify header
			file, err := os.Open(expectedFileName)
			if err != nil {
				t.Fatalf("Failed to open CSV file for verification: %v", err)
			}
			defer file.Close()

			reader := csv.NewReader(file)
			header, err := reader.Read()
			if err != nil {
				t.Fatalf("Failed to read header from CSV file: %v", err)
			}

			// Verify header matches expected
			if len(header) != len(tt.expectedHeader) {
				t.Errorf("Header length = %d, want %d. Got: %v, Want: %v",
					len(header), len(tt.expectedHeader), header, tt.expectedHeader)
				return
			}

			for i, col := range header {
				if col != tt.expectedHeader[i] {
					t.Errorf("Header column %d = %q, want %q", i, col, tt.expectedHeader[i])
				}
			}

			// Cleanup - remove the test file
			os.Remove(expectedFileName)
		})
	}
}

// TestCSVLogger_WriteRow tests writing rows to the CSV logger
func TestCSVLogger_WriteRow(t *testing.T) {
	// Create a test logger
	testAction := "test_writerow"
	logger, err := NewCSVLogger(testAction)
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}
	defer func() {
		logger.Close()
		// Cleanup
		tempDir := os.TempDir()
		dateStr := time.Now().Format("2006-01-02")
		testFile := filepath.Join(tempDir, "_msgraphtool_"+testAction+"_"+dateStr+".csv")
		os.Remove(testFile)
	}()

	// Note: Default action header has 4 columns: Timestamp, Action, Status, Details
	// We need to provide 3 columns (Timestamp is added automatically by WriteRow)
	tests := []struct {
		name string
		row  []string
	}{
		{
			name: "standard row",
			row:  []string{testAction, "Success", "Test details"},
		},
		{
			name: "row with empty values",
			row:  []string{testAction, "", ""},
		},
		{
			name: "row with special characters",
			row:  []string{testAction, "Success", "Details with, comma and \"quotes\""},
		},
		{
			name: "row with long text",
			row:  []string{testAction, "Success", "This is a longer detail field with multiple words and data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("WriteRow() panicked: %v", r)
				}
			}()

			logger.WriteRow(tt.row)
		})
	}

	// Close to flush
	if err := logger.Close(); err != nil {
		t.Fatalf("Failed to close logger: %v", err)
	}

	// Verify rows were written
	tempDir := os.TempDir()
	dateStr := time.Now().Format("2006-01-02")
	testFile := filepath.Join(tempDir, "_msgraphtool_"+testAction+"_"+dateStr+".csv")

	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file for verification: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	// Should have header + number of test rows
	expectedRows := 1 + len(tests) // header + test rows
	if len(records) != expectedRows {
		t.Errorf("CSV file has %d rows, want %d (1 header + %d data rows)",
			len(records), expectedRows, len(tests))
	}

	// Verify each data row has a timestamp prepended
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 1 {
			t.Errorf("Row %d is empty", i)
			continue
		}

		// First column should be timestamp
		timestamp := records[i][0]
		if timestamp == "" {
			t.Errorf("Row %d missing timestamp", i)
		}

		// Verify timestamp format (should be parseable)
		_, err := time.Parse("2006-01-02 15:04:05", timestamp)
		if err != nil {
			t.Errorf("Row %d has invalid timestamp format: %q, error: %v", i, timestamp, err)
		}
	}
}

// TestCSVLogger_Close tests closing the CSV logger properly
func TestCSVLogger_Close(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*CSVLogger)
		wantErr   bool
	}{
		{
			name: "close empty logger",
			setupFunc: func(l *CSVLogger) {
				// No writes
			},
			wantErr: false,
		},
		{
			name: "close after single write",
			setupFunc: func(l *CSVLogger) {
				l.WriteRow([]string{"test", "Success", "Detail"})
			},
			wantErr: false,
		},
		{
			name: "close after multiple writes",
			setupFunc: func(l *CSVLogger) {
				for i := 0; i < 5; i++ {
					l.WriteRow([]string{"test", "Success", "Row data"})
				}
			},
			wantErr: false,
		},
		{
			name: "double close returns error",
			setupFunc: func(l *CSVLogger) {
				l.WriteRow([]string{"test", "Success", "Data"})
				l.Close() // First close
			},
			wantErr: true, // Second close returns error (file already closed)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAction := "test_close_" + strings.ReplaceAll(tt.name, " ", "_")
			logger, err := NewCSVLogger(testAction)
			if err != nil {
				t.Fatalf("Failed to create test logger: %v", err)
			}

			// Cleanup file
			defer func() {
				tempDir := os.TempDir()
				dateStr := time.Now().Format("2006-01-02")
				testFile := filepath.Join(tempDir, "_msgraphtool_"+testAction+"_"+dateStr+".csv")
				os.Remove(testFile)
			}()

			// Setup
			tt.setupFunc(logger)

			// Test close
			err = logger.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify file is properly closed by trying to read it
			tempDir := os.TempDir()
			dateStr := time.Now().Format("2006-01-02")
			testFile := filepath.Join(tempDir, "_msgraphtool_"+testAction+"_"+dateStr+".csv")

			// Should be able to open and read the file after close
			file, err := os.Open(testFile)
			if err != nil {
				t.Errorf("Failed to open CSV file after close: %v", err)
				return
			}
			defer file.Close()

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				t.Errorf("Failed to read CSV file after close: %v", err)
			}

			// Should have at least a header
			if len(records) < 1 {
				t.Error("CSV file has no header after close")
			}
		})
	}
}

// TestCSVLogger_MultipleWrites tests writing multiple rows and flushing behavior
func TestCSVLogger_MultipleWrites(t *testing.T) {
	testAction := "test_multiwrites"
	logger, err := NewCSVLogger(testAction)
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}

	// Cleanup
	defer func() {
		logger.Close()
		tempDir := os.TempDir()
		dateStr := time.Now().Format("2006-01-02")
		testFile := filepath.Join(tempDir, "_msgraphtool_"+testAction+"_"+dateStr+".csv")
		os.Remove(testFile)
	}()

	// Test writing more than flushEvery rows to trigger auto-flush
	const numRows = 25 // flushEvery is 10, so this should trigger 2 auto-flushes

	// Note: Default action header has 4 columns: Timestamp, Action, Status, Details
	// WriteRow prepends timestamp, so we provide 3 columns
	for i := 0; i < numRows; i++ {
		logger.WriteRow([]string{testAction, "Success", "Row data " + string(rune('A'+(i%26)))})
	}

	// Verify rowCount is tracked correctly
	if logger.rowCount != numRows {
		t.Errorf("Logger rowCount = %d, want %d", logger.rowCount, numRows)
	}

	// Close to flush remaining data
	if err := logger.Close(); err != nil {
		t.Fatalf("Failed to close logger: %v", err)
	}

	// Read and verify all rows were written
	tempDir := os.TempDir()
	dateStr := time.Now().Format("2006-01-02")
	testFile := filepath.Join(tempDir, "_msgraphtool_"+testAction+"_"+dateStr+".csv")

	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file for verification: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	// Should have header + numRows data rows
	expectedRecords := 1 + numRows
	if len(records) != expectedRecords {
		t.Errorf("CSV file has %d rows, want %d (1 header + %d data rows)",
			len(records), expectedRecords, numRows)
	}

	// Verify all rows have timestamps
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 1 {
			t.Errorf("Row %d is empty", i)
			continue
		}

		// Verify timestamp exists and is valid
		timestamp := records[i][0]
		if _, err := time.Parse("2006-01-02 15:04:05", timestamp); err != nil {
			t.Errorf("Row %d has invalid timestamp: %q", i, timestamp)
		}
	}
}

// TestCSVLogger_AppendMode tests that CSV logger appends to existing files
func TestCSVLogger_AppendMode(t *testing.T) {
	testAction := "test_append"
	tempDir := os.TempDir()
	dateStr := time.Now().Format("2006-01-02")
	testFile := filepath.Join(tempDir, "_msgraphtool_"+testAction+"_"+dateStr+".csv")

	// Cleanup
	defer os.Remove(testFile)

	// First logger - create file and write some rows
	logger1, err := NewCSVLogger(testAction)
	if err != nil {
		t.Fatalf("Failed to create first logger: %v", err)
	}

	logger1.WriteRow([]string{testAction, "Success", "First logger data"})
	logger1.WriteRow([]string{testAction, "Success", "First logger data 2"})

	if err := logger1.Close(); err != nil {
		t.Fatalf("Failed to close first logger: %v", err)
	}

	// Read file to get row count after first logger
	file1, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open file after first logger: %v", err)
	}
	reader1 := csv.NewReader(file1)
	records1, err := reader1.ReadAll()
	file1.Close()
	if err != nil {
		t.Fatalf("Failed to read file after first logger: %v", err)
	}

	firstLoggerRows := len(records1)
	t.Logf("After first logger: %d rows (1 header + 2 data rows)", firstLoggerRows)

	// Second logger - should append to existing file
	logger2, err := NewCSVLogger(testAction)
	if err != nil {
		t.Fatalf("Failed to create second logger: %v", err)
	}

	logger2.WriteRow([]string{testAction, "Success", "Second logger data"})
	logger2.WriteRow([]string{testAction, "Success", "Second logger data 2"})

	if err := logger2.Close(); err != nil {
		t.Fatalf("Failed to close second logger: %v", err)
	}

	// Read file to verify append
	file2, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open file after second logger: %v", err)
	}
	defer file2.Close()

	reader2 := csv.NewReader(file2)
	records2, err := reader2.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read file after second logger: %v", err)
	}

	// Should have: 1 header (from first logger) + 2 rows (first logger) + 2 rows (second logger) = 5 total
	expectedTotalRows := 5
	if len(records2) != expectedTotalRows {
		t.Errorf("After append, file has %d rows, want %d rows", len(records2), expectedTotalRows)
	}

	// Verify header appears only once (at the beginning)
	if records2[0][0] != "Timestamp" {
		t.Error("First row should be header with Timestamp column")
	}

	// Verify subsequent rows have timestamps (not headers)
	for i := 1; i < len(records2); i++ {
		if records2[i][0] == "Timestamp" {
			t.Errorf("Row %d is a duplicate header, should be data row", i)
		}

		// Verify it's a valid timestamp
		if _, err := time.Parse("2006-01-02 15:04:05", records2[i][0]); err != nil {
			t.Errorf("Row %d has invalid timestamp: %q", i, records2[i][0])
		}
	}

	t.Logf("✓ Append mode verified: %d total rows (%d from first logger + %d from second logger)",
		len(records2), firstLoggerRows-1, len(records2)-firstLoggerRows)
}

// TestCSVLogger_NilWriterHandling tests that WriteRow handles nil writer gracefully
func TestCSVLogger_NilWriterHandling(t *testing.T) {
	// Create a logger with nil writer (simulating a closed logger)
	logger := &CSVLogger{
		writer:     nil,
		file:       nil,
		action:     "test",
		rowCount:   0,
		lastFlush:  time.Now(),
		flushEvery: 10,
	}

	// WriteRow should not panic with nil writer
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("WriteRow() with nil writer panicked: %v", r)
		}
	}()

	logger.WriteRow([]string{"test", "data"})

	// Close should not panic with nil writer/file
	if err := logger.Close(); err != nil {
		t.Errorf("Close() with nil writer/file returned error: %v", err)
	}
}
