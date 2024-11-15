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

	// if the logger is enabled, write a log entry before value is created/updated
	kv.quickLog("SET", key, string(value))

	// Set value in store
	kv.data[key] = value
	log.Printf("Set \"%s\" to %s", key, value)
}

// Get retrieves the value associated with a key from the store.
func (kv *KeyValueStore) Get(key string) (json.RawMessage, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	val, ok := kv.data[key]
	log.Printf("Get \"%s\" from kvs", key)

	// Write log entry
	kv.quickLog("GET", key, string(val))

	return val, ok
}

// GetAll retries all key-values pairs from the store.
func (kv *KeyValueStore) GetAll() map[string][]byte {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	// Log operation, even though we're not logging the values themselves
	kv.quickLog("GETALL", "", "")
	log.Print("Get all key-value pairs from kvs")

	return maps.Clone(kv.data)
}

// GetKeys returns all keys from the store.
func (kv *KeyValueStore) GetKeys() []string {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	// Log operation, even though we're not logging the keys themselves
	kv.quickLog("GETKEYS", "", "")
	log.Print("Get all keys from kvs")

	// Copy keys to a new slice
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

	// Log operation, even though we're not logging the values themselves
	kv.quickLog("GETVALUES", "", "")
	log.Print("Get all values from kvs")

	// Copy values to a new slice
	values := make([]json.RawMessage, 0, len(kv.data))
	for _, v := range kv.data {
		values = append(values, v)
	}

	return values
}

// Deletes all key/value pairs from the store and clears the transaction logs.
func (kv *KeyValueStore) DeleteALL() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Log the operation
	kv.quickLog("DELETEALL", "", "")

	// Clear the in-memory data
	kv.data = make(map[string][]byte)
	log.Print("Delete all key-value pairs from kvs")

	return nil
}

// Deletes a specific key value pair from the store.
func (kv *KeyValueStore) Delete(key string) ([]byte, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// get value to be deleted
	deletedVal, ok := kv.data[key]

	// log entry
	kv.quickLog("DELETE", key, string(deletedVal))

	// delete key from store
	delete(kv.data, key)
	log.Printf("Deleted \"%s\" from kvs", key)

	return deletedVal, ok
}

// ProcessLogEntries processes a list of log entries and updates the key-value store accordingly.
func (kv *KeyValueStore) ProcessLogEntries(entries []LogEntry) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Process each log entry depending on the operation
	for _, entry := range entries {
		switch entry.Operation {
		case "SET": // Add or update the key:value pair in the in-memory data
			kv.data[entry.Key] = []byte(entry.Value)
		case "DELETE": // Delete the key:value pair from the in-memory data
			delete(kv.data, entry.Key)
		case "DELETEALL": // Clear all the in-memory data
			kv.data = make(map[string][]byte)
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

// Quick log entry utility function.
func (kv *KeyValueStore) quickLog(operation string, key string, value string) {
	if kv.logger != nil {
		kv.logger.WriteLog(LogEntry{
			Timestamp: time.Now(),
			Operation: operation,
			Key:       key,
			Value:     value,
		})
	}
}
