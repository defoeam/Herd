package keyvaluestore

import (
	"fmt"
	"log"
	"maps"
	"os"
	"sync"
	"time"

	"github.com/goccy/go-json"
)

// KeyValueStore represents the key-value store.
type KeyValueStore struct {
	data   map[string][]byte
	mu     sync.RWMutex
	logger *Logger
}

// NewKeyValueStore creates a new instance of KeyValueStore.
func NewKeyValueStore(logFile string) (*KeyValueStore, error) {
	// Check if the log file exists, create it if it doesn't
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, err := os.Create(logFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %v", err)
		}
		file.Close()
	}

	logger, err := NewLogger(logFile)
	if err != nil {
		return nil, err
	}

	kv := &KeyValueStore{
		data:   make(map[string][]byte),
		logger: logger,
	}

	// Read and process log entries
	entries, err := logger.ReadLogs()
	if err != nil {
		return nil, fmt.Errorf("failed to read log entries: %v", err)
	}

	kv.ProcessLogEntries(entries)

	return kv, nil
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

	kv.logger.Log(LogEntry{
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

// Clears all key/value pairs from the store.
func (kv *KeyValueStore) ClearAll() {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	for k := range kv.data {
		delete(kv.data, k)
	}
}

// Clears a specific key value pair from the store.
func (kv *KeyValueStore) Clear(key string) ([]byte, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	deletedVal, ok := kv.data[key]
	delete(kv.data, key)

	kv.logger.Log(LogEntry{
		Timestamp: time.Now(),
		Operation: "DELETE",
		Key:       key,
		Value:     string(deletedVal),
	})

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

func (kv *KeyValueStore) CompactLogs() error {
	return kv.logger.CompactLogs()
}
