package keyvaluestore

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// LogEntry represents a log entry.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
}

// Logger is a simple logger that writes to a file.
type Logger struct {
	filename string
	mu       sync.RWMutex
}

// NewLogger creates a new logger that writes to the specified file.
func NewLogger(filename string) (*Logger, error) {
	// Attempt to open the file to ensure it exists and is accessible
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open or create log file: %w", err)
	}
	file.Close() // Close the file as we don't need to keep it open

	return &Logger{filename: filename}, nil
}

// WriteLog writes a log entry to the logger's file.
func (l *Logger) WriteLog(entry LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Open the log file for appending
	file, err := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	// Format the log entry
	logLine := fmt.Sprintf("[%s] %s - Key: %s, Value: %s\n",
		entry.Timestamp.Format(time.RFC3339),
		entry.Operation,
		entry.Key,
		entry.Value,
	)

	if _, writeErr := file.WriteString(logLine); writeErr != nil {
		log.Printf("Error writing to log file: %v", writeErr)
	}
}

// parseLogLine parses a log line into a LogEntry struct.
func parseLogLine(line string) (LogEntry, error) {
	const (
		splitParts        = 2
		keyValueSeparator = ", "
	)

	// Split the log line into timestamp and operation/key-value parts
	parts := strings.SplitN(line, "] ", splitParts)
	if len(parts) != splitParts {
		return LogEntry{}, errors.New("invalid log line format")
	}

	// Parse the timestamp
	timestamp, err := time.Parse(time.RFC3339, strings.Trim(parts[0], "[]"))
	if err != nil {
		return LogEntry{}, err
	}

	// Parse the operation and key-value parts
	operationParts := strings.SplitN(parts[1], " - ", splitParts)
	if len(operationParts) != splitParts {
		return LogEntry{}, errors.New("invalid log line format (operation)")
	}

	// Parse the key and value
	operation := operationParts[0]
	keyValue := strings.SplitN(operationParts[1], keyValueSeparator, splitParts)
	if len(keyValue) != splitParts {
		return LogEntry{}, errors.New("invalid log line format (key/value)")
	}

	key := strings.TrimPrefix(keyValue[0], "Key: ")
	value := strings.TrimPrefix(keyValue[1], "Value: ")

	return LogEntry{
		Timestamp: timestamp,
		Operation: operation,
		Key:       key,
		Value:     value,
	}, nil
}

// ReadLogs reads all log entries from the file.
func (l *Logger) ReadLogs() ([]LogEntry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Open the log file for reading
	file, err := os.Open(l.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read each line from the file and parse it into a LogEntry
	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		entry, parseLogLineErr := parseLogLine(line)
		if parseLogLineErr != nil {
			log.Printf("Error parsing log line: %v", parseLogLineErr)
			continue
		}
		entries = append(entries, entry)
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, scanErr
	}

	return entries, nil
}

// ClearLogs clears the entire transaction log file.
func (l *Logger) ClearLogs() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Truncate the file to zero length
	err := os.Truncate(l.filename, 0)
	if err != nil {
		return fmt.Errorf("failed to truncate log file: %w", err)
	}

	return nil
}
