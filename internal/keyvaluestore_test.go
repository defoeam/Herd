package keyvaluestore_test

import (
	"encoding/json"
	"log"
	"testing"

	herd "github.com/defoeam/herd/internal"
)

func TestGetSet(t *testing.T) {
	// Create a new KeyValueStore instance
	kv := herd.NewKeyValueStore()

	t.Run("Set and Get", func(t *testing.T) {
		key := "test_key"
		value := json.RawMessage(`{"name":"John","age":30}`)

		kv.Set(key, value)

		retrievedValue, ok := kv.Get(key)
		if !ok {
			t.Errorf("Failed to retrieve value for key %s", key)
		}

		if string(retrievedValue) != string(value) {
			t.Errorf("Retrieved value %s does not match set value %s", retrievedValue, value)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		kv.Set("key1", json.RawMessage(`"value1"`))
		kv.Set("key2", json.RawMessage(`"value2"`))

		allItems := kv.GetAll()

		if len(allItems) != 3 { // Including the previous "test_key"
			t.Errorf("Expected 3 items, got %d", len(allItems))
		}

		if string(allItems["key1"]) != `"value1"` {
			t.Errorf("Unexpected value for key1: %s", allItems["key1"])
		}

		if string(allItems["key2"]) != `"value2"` {
			t.Errorf("Unexpected value for key2: %s", allItems["key2"])
		}
	})
}

func TestGetKeysValues(t *testing.T) {
	kv := herd.NewKeyValueStore()

	t.Run("GetKeys", func(t *testing.T) {
		// clear all items to prep for test
		err := kv.ClearAll()
		if err != nil {
			log.Printf("Failed to clear all items: %v", err)
		}

		// add test key/value pairs
		kv.Set("key1", json.RawMessage(`"value1"`))
		kv.Set("key2", json.RawMessage(`"value2"`))

		keys := kv.GetKeys()

		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d", len(keys))
		}

		expectedKeys := map[string]bool{"key1": true, "key2": true}
		for _, key := range keys {
			if !expectedKeys[key] {
				t.Errorf("Unexpected key: %s", key)
			}
		}
	})

	t.Run("GetValues", func(t *testing.T) {
		// clear all items to prep for test
		err := kv.ClearAll()
		if err != nil {
			log.Printf("Failed to clear all items: %v", err)
		}

		// add test key/value pairs
		kv.Set("key1", json.RawMessage(`"value1"`))
		kv.Set("key2", json.RawMessage(`"value2"`))

		values := kv.GetValues()

		if len(values) != 2 {
			t.Errorf("Expected 2 values, got %d", len(values))
		}

		expectedValues := map[string]bool{`"value1"`: true, `"value2"`: true}
		for _, value := range values {
			if !expectedValues[string(value)] {
				t.Errorf("Unexpected value: %s", string(value))
			}
		}
	})
}

func TestClear(t *testing.T) {
	// Create a new KeyValueStore instance
	kv := herd.NewKeyValueStore()

	t.Run("Clear", func(t *testing.T) {
		key := "key_to_clear"
		value := json.RawMessage(`"clear_me"`)

		kv.Set(key, value)

		clearedValue, ok := kv.Clear(key)
		if !ok {
			t.Errorf("Failed to clear key %s", key)
		}

		if string(clearedValue) != string(value) {
			t.Errorf("Cleared value %s does not match set value %s", clearedValue, value)
		}

		_, exists := kv.Get(key)
		if exists {
			t.Errorf("Key %s still exists after clearing", key)
		}
	})

	t.Run("ClearAll", func(t *testing.T) {
		kv.Set("key1", json.RawMessage(`"value1"`))
		kv.Set("key2", json.RawMessage(`"value2"`))

		if err := kv.ClearAll(); err != nil {
			// Handle the error appropriately
			log.Printf("Failed to clear all items: %v", err)
		}

		allItems := kv.GetAll()
		if len(allItems) != 0 {
			t.Errorf("Expected 0 items after ClearAll, got %d", len(allItems))
		}
	})
}
