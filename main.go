package main

import (
	"fmt"
	"log"
	"net/http"

	keyvaluestore "github.com/defoeam/kvs/kv"
	"github.com/gin-gonic/gin"
)

func handleGet(kv *keyvaluestore.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Param("key")

		// Check if empty key
		if key == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Key parameter is missing"})
			return
		}

		// Check if key exists and get value
		val, ok := kv.Get(key)
		if !ok {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Key not found"})
			return
		}

		// Build response
		res := struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{Key: key, Value: val}

		// Serialize response
		ctx.JSON(http.StatusAccepted, res)
	}
}

func handleSet(kv *keyvaluestore.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Define key/value request structure
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}

		// Bind JSON to key/value structure
		err := ctx.BindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind json, incorrect parameter(s)"})
			return
		}

		// Add to kv storage
		kv.Set(req.Key, req.Value)

		// Serialize response
		ctx.JSON(http.StatusCreated, req)
	}
}

func main() {

	// Create a new instance of KeyValueStore.
	kv := keyvaluestore.NewKeyValueStore()

	// Setup gin engine
	router := gin.Default()

	// GET endpoints
	router.GET("/get/:key", handleGet(kv))

	// POST endpoints
	router.POST("/set", handleSet(kv))

	// Get address
	port := 8080
	addr := fmt.Sprintf(":%d", port)

	log.Printf("Starting key-value store on http://localhost%s\n", addr)
	log.Fatal(router.Run(addr))

}
