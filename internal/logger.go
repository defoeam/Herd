package keyvaluestore

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		return nil, fmt.Errorf("failed to open or create log file: %v", err)
	}
	file.Close() // Close the file as we don't need to keep it open

	return &Logger{filename: filename}, nil
}

// Log writes a log entry to the logger's file.
func (l *Logger) Log(entry LogEntry) {
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

	if _, err := file.WriteString(logLine); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

// ReadLogs reads all log entries from the file.
func (l *Logger) ReadLogs() ([]LogEntry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.readLogsUnsafe()
}

func parseLogLine(line string) (LogEntry, error) {
	// Example log line format:
	// [2023-04-13T12:34:56Z] SET - Key: mykey, Value: {"foo": "bar"}
	parts := strings.SplitN(line, "] ", 2)
	if len(parts) != 2 {
		return LogEntry{}, fmt.Errorf("invalid log line format")
	}

	timestamp, err := time.Parse(time.RFC3339, strings.Trim(parts[0], "[]"))
	if err != nil {
		return LogEntry{}, err
	}

	operationParts := strings.SplitN(parts[1], " - ", 2)
	if len(operationParts) != 2 {
		return LogEntry{}, fmt.Errorf("invalid log line format")
	}

	operation := operationParts[0]
	keyValue := strings.SplitN(operationParts[1], ", ", 2)
	if len(keyValue) != 2 {
		return LogEntry{}, fmt.Errorf("invalid log line format")
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

func (l *Logger) CompactLogs() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Read all log entries
	entries, err := l.readLogsUnsafe()
	if err != nil {
		return fmt.Errorf("failed to read log entries: %v", err)
	}

	// Keep only the latest SET entry for each key, and track DELETEs
	latestEntries := make(map[string]LogEntry)
	deletedKeys := make(map[string]bool)
	for _, entry := range entries {
		switch entry.Operation {
		case "SET":
			latestEntries[entry.Key] = entry
			delete(deletedKeys, entry.Key)
		case "DELETE":
			delete(latestEntries, entry.Key)
			deletedKeys[entry.Key] = true
		}
	}

	// Create a temporary file for writing compacted logs
	tempFile, err := os.CreateTemp(filepath.Dir(l.filename), "compacted_*.log")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tempFile.Close()

	// Write compacted logs to the temporary file
	for _, entry := range latestEntries {
		logLine := fmt.Sprintf("[%s] %s - Key: %s, Value: %s\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.Operation,
			entry.Key,
			entry.Value,
		)
		if _, err := tempFile.WriteString(logLine); err != nil {
			return fmt.Errorf("failed to write to temporary file: %v", err)
		}
	}

	// Write DELETE entries for keys that were ultimately deleted
	for key := range deletedKeys {
		logLine := fmt.Sprintf("[%s] DELETE - Key: %s, Value: \n",
			time.Now().Format(time.RFC3339),
			key,
		)
		if _, err := tempFile.WriteString(logLine); err != nil {
			return fmt.Errorf("failed to write to temporary file: %v", err)
		}
	}

	// Close the temporary file
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %v", err)
	}

	// Rename the temporary file to replace the original log file
	if err := os.Rename(tempFile.Name(), l.filename); err != nil {
		return fmt.Errorf("failed to rename temporary file: %v", err)
	}

	return nil
}

// readLogsUnsafe reads log entries without locking the mutex
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
		entry, err := parseLogLine(line)
		if err != nil {
			log.Printf("Error parsing log line: %v", err)
			continue
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
