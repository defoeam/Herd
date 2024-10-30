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
	Timestamp time.Time
	Operation string
	Key       string
	Value     string
}

// Logger is a simple logger that writes to a file.
type Logger struct {
	filename string
	mu       sync.Mutex
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

	file, err := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

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

// ReadLogs reads all log entries from the file.
func (l *Logger) ReadLogs() ([]LogEntry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.readLogsUnsafe()
}

func parseLogLine(line string) (LogEntry, error) {
	const (
		splitParts        = 2
		keyValueSeparator = ", "
	)

	parts := strings.SplitN(line, "] ", splitParts)
	if len(parts) != splitParts {
		return LogEntry{}, errors.New("invalid log line format")
	}

	timestamp, err := time.Parse(time.RFC3339, strings.Trim(parts[0], "[]"))
	if err != nil {
		return LogEntry{}, err
	}

	operationParts := strings.SplitN(parts[1], " - ", splitParts)
	if len(operationParts) != splitParts {
		return LogEntry{}, errors.New("invalid log line format (operation)")
	}

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

// readLogsUnsafe reads log entries without locking the mutex.
func (l *Logger) readLogsUnsafe() ([]LogEntry, error) {
	file, err := os.Open(l.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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
