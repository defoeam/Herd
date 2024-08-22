package port

import "encoding/json"

type StoreRepository interface {
	Set(key string, value json.RawMessage)
	Get(key string) (json.RawMessage, bool)
	GetAll() map[string][]byte
	GetKeys() []string
	GetValues() []json.RawMessage
	ClearAll()
	Clear(key string) ([]byte, bool)
}
