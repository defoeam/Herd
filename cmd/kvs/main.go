package main

import (
	"fmt"
	"log"

	"github.com/defoeam/kvs/internal/adapter/handler/http"
	"github.com/defoeam/kvs/internal/adapter/keyvaluestore"
	"github.com/defoeam/kvs/internal/adapter/keyvaluestore/keyvaluestorememory"
)

func main() {
	var memoryStore = keyvaluestorememory.NewKeyValueMemoryStore()
	var keyvaluestore = keyvaluestore.NewKeyValueStore(memoryStore)

	var router = http.NewHttp(keyvaluestore)

	// Get address
	port := 8080
	addr := fmt.Sprintf(":%d", port)

	// Start server
	log.Printf("Starting key-value store on http://localhost%s\n", addr)
	log.Fatal(router.Run(addr))
}
