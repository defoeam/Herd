package keyvaluestorememory

import (
	"encoding/json"
	"log"
	"maps"
	"sync"

	"github.com/defoeam/kvs/internal/core/port"
)

// KeyValueMemoryStore represents the key-value store.
type KeyValueMemoryStore struct {
	data map[string][]byte
	mu   sync.RWMutex
}

// NewKeyValueMemoryStore creates a new instance of KeyValueMemoryStore.
func NewKeyValueMemoryStore() port.StoreRepository {
	return &KeyValueMemoryStore{
		data: make(map[string][]byte),
	}
}

// Set adds or updates a key-value pair in the store.
func (kv *KeyValueMemoryStore) Set(key string, value json.RawMessage) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	log.Printf("Adding \"%s\" to kvs", key)
	kv.data[key] = value
}

// Get retrieves the value associated with a key from the store.
func (kv *KeyValueMemoryStore) Get(key string) (json.RawMessage, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	val, ok := kv.data[key]
	return val, ok
}

// GetAll retries all key-values pairs from the store.
func (kv *KeyValueMemoryStore) GetAll() map[string][]byte {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return maps.Clone(kv.data)
}

// GetKeys returns all keys from the store.
func (kv *KeyValueMemoryStore) GetKeys() []string {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	keys := make([]string, 0, len(kv.data))
	for k := range kv.data {
		keys = append(keys, k)
	}

	return keys
}

// GetValues returns all values from the store.
func (kv *KeyValueMemoryStore) GetValues() []json.RawMessage {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	values := make([]json.RawMessage, 0, len(kv.data))
	for _, v := range kv.data {
		values = append(values, v)
	}

	return values
}

// Clears all key/value pairs from the store.
func (kv *KeyValueMemoryStore) ClearAll() {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	for k := range kv.data {
		delete(kv.data, k)
	}
}

// Clears a specific key value pair from the store.
func (kv *KeyValueMemoryStore) Clear(key string) ([]byte, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	deletedVal, ok := kv.data[key]
	delete(kv.data, key)

	return deletedVal, ok
}
