package keyvaluestore

import (
	"log"
	"maps"
	"sync"
)

// KeyValueStore represents the key-value store.
type KeyValueStore struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewKeyValueStore creates a new instance of KeyValueStore.
func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		data: make(map[string]string),
	}
}

// Set adds or updates a key-value pair in the store.
func (kv *KeyValueStore) Set(key string, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	log.Printf("Adding \"%s\" to kvs", key)
	kv.data[key] = value
}

// Get retrieves the value associated with a key from the store.
func (kv *KeyValueStore) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	val, ok := kv.data[key]
	return val, ok
}

// GetAll retries all key-values pairs from the store.
func (kv *KeyValueStore) GetAll() map[string]string {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return maps.Clone(kv.data)
}

func (kv *KeyValueStore) GetKeys() []string {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	keys := make([]string, 0, len(kv.data))
	for k := range kv.data {
		keys = append(keys, k)
	}

	return keys
}

func (kv *KeyValueStore) GetValues() []string {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	values := make([]string, 0, len(kv.data))
	for _, v := range kv.data {
		values = append(values, v)
	}

	return values
}

func (kv *KeyValueStore) Clear() {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	for k := range kv.data {
		delete(kv.data, k)
	}
}
