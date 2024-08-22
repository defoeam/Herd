package keyvaluestore

import (
	"encoding/json"

	"github.com/defoeam/kvs/internal/core/port"
)

type KeyValueStore struct {
	store port.StoreRepository
}

func NewKeyValueStore(store port.StoreRepository) port.KeyValueStore {
	return &KeyValueStore{
		store: store,
	}
}

func (kvs KeyValueStore) Set(key string, value json.RawMessage) {
	kvs.store.Set(key, value)
}
func (kvs KeyValueStore) Get(key string) (json.RawMessage, bool) {
	return kvs.store.Get(key)
}
func (kvs KeyValueStore) GetAll() map[string][]byte {
	return kvs.store.GetAll()
}
func (kvs KeyValueStore) GetKeys() []string {
	return kvs.store.GetKeys()
}
func (kvs KeyValueStore) GetValues() []json.RawMessage {
	return kvs.store.GetValues()
}
func (kvs KeyValueStore) ClearAll() {
	kvs.store.ClearAll()
}
func (kvs KeyValueStore) Clear(key string) ([]byte, bool) {
	return kvs.store.Clear(key)
}
