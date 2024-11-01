package keyvaluestore

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"sync"
	"time"
)

// KeyValueStore represents the key-value store.
type KeyValueStore struct {
	data             map[string][]byte
	mu               sync.RWMutex
	logger           *Logger
	snapshotInterval time.Duration
}

func (kv *KeyValueStore) InitLogging(logFile string, snapshotInterval time.Duration) error {
	// Check if the log file exists, create it if it doesn't
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, createErr := os.Create(logFile)
		if createErr != nil {
			return fmt.Errorf("failed to create log file: %w", createErr)
		}
		file.Close()
	}

	logger, loggerErr := NewLogger(logFile)
	if loggerErr != nil {
		return loggerErr
	}

	kv.logger = logger
	kv.snapshotInterval = snapshotInterval

	// Load the latest snapshot
	if snapshotErr := kv.LoadLatestSnapshot(); snapshotErr != nil {
		return fmt.Errorf("failed to load latest snapshot: %w", snapshotErr)
	}

	// Read and process log entries
	entries, readLogsErr := logger.ReadLogs()
	if readLogsErr != nil {
		return fmt.Errorf("failed to read log entries: %w", readLogsErr)
	}

	kv.ProcessLogEntries(entries)

	// Start the snapshot scheduler
	go kv.snapshotScheduler()
	return nil
}

// NewKeyValueStore creates a new instance of KeyValueStore.
func NewKeyValueStore() *KeyValueStore {
	kv := &KeyValueStore{
		data:             make(map[string][]byte),
		logger:           nil,
		snapshotInterval: 1 * time.Hour,
	}

	return kv
}

type KeyValue struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

func NewKeyValue(key string, value json.RawMessage) *KeyValue {
	kv := KeyValue{
		Key:   key,
		Value: value,
	}

	return &kv
}

// Set adds or updates a key-value pair in the store.
func (kv *KeyValueStore) Set(key string, value json.RawMessage) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	log.Printf("Adding \"%s\" to kvs", key)
	kv.data[key] = value

	// if the logger is enabled, write a log entry
	if kv.logger == nil {
		return
	}

	kv.logger.WriteLog(LogEntry{
		Timestamp: time.Now(),
		Operation: "SET",
		Key:       key,
		Value:     string(value),
	})
}

// Get retrieves the value associated with a key from the store.
func (kv *KeyValueStore) Get(key string) (json.RawMessage, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	val, ok := kv.data[key]
	return val, ok
}

// GetAll retries all key-values pairs from the store.
func (kv *KeyValueStore) GetAll() map[string][]byte {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return maps.Clone(kv.data)
}

// GetKeys returns all keys from the store.
func (kv *KeyValueStore) GetKeys() []string {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	keys := make([]string, 0, len(kv.data))
	for k := range kv.data {
		keys = append(keys, k)
	}

	return keys
}

// GetValues returns all values from the store.
func (kv *KeyValueStore) GetValues() []json.RawMessage {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	values := make([]json.RawMessage, 0, len(kv.data))
	for _, v := range kv.data {
		values = append(values, v)
	}

	return values
}

// Clears all key/value pairs from the store and clears the transaction logs.
func (kv *KeyValueStore) ClearAll() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Clear the in-memory data
	kv.data = make(map[string][]byte)

	// if the logger is enabled, write a log entry
	if kv.logger == nil {
		return nil
	}

	// Clear the transaction logs
	err := kv.logger.ClearLogs()
	if err != nil {
		return fmt.Errorf("failed to clear transaction logs: %w", err)
	}

	return nil
}

// Clears a specific key value pair from the store.
func (kv *KeyValueStore) Clear(key string) ([]byte, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	deletedVal, ok := kv.data[key]
	delete(kv.data, key)

	if kv.logger != nil {
		kv.logger.WriteLog(LogEntry{
			Timestamp: time.Now(),
			Operation: "DELETE",
			Key:       key,
			Value:     string(deletedVal),
		})
	}

	return deletedVal, ok
}

func (kv *KeyValueStore) ProcessLogEntries(entries []LogEntry) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	for _, entry := range entries {
		switch entry.Operation {
		case "SET":
			kv.data[entry.Key] = []byte(entry.Value)
		case "DELETE":
			delete(kv.data, entry.Key)
		}
	}
}

// snapshotScheduler runs periodically to take snapshots of the key-value store.
// It uses a ticker to trigger snapshots at the interval specified by kv.snapshotInterval.
// If a snapshot fails, it logs the error but continues running.
func (kv *KeyValueStore) snapshotScheduler() {
	ticker := time.NewTicker(kv.snapshotInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := kv.TakeSnapshot(); err != nil {
			log.Printf("Failed to take snapshot: %v", err)
		}

		log.Printf("Snapshot taken at %v", time.Now())
	}
}
