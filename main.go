package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	keyvaluestore "github.com/defoeam/kvs/kv"
)

// HandleSet handles the HTTP endpoint for setting key-value pairs.
func HandleSet(kv *keyvaluestore.KeyValueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Key   string `json:"key"`
			Value string `json:"Value"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		kv.Set(req.Key, req.Value)
		w.WriteHeader(http.StatusOK)
	}
}

// HandleGet handles the HTTP endpoint for retrieving value by key.
func HandleGet(kv *keyvaluestore.KeyValueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "Key parameter is missing", http.StatusBadRequest)
			return
		}

		val, ok := kv.Get(key)
		if !ok {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}

		resp := struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{Key: key, Value: val}

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "Failed to json encode", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// Create a new instance of KeyValueStore.
	kv := keyvaluestore.NewKeyValueStore()
	// Set up HTTP handlers
	http.HandleFunc("/set", HandleSet(kv))
	http.HandleFunc("/get", HandleGet(kv))

	// Timeout multiplier, this fixes the magic number linting error
	timeoutMultiplier := 10

	// Start the HTTP server.
	port := 8080
	addr := fmt.Sprintf(":%d", port)
	readTimeout := time.Duration(timeoutMultiplier) * time.Second
	writeTimeout := time.Duration(timeoutMultiplier) * time.Second
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Printf("Starting key-value store on http://localhost%s\n", addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
