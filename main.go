package main

import (
	"fmt"
	"log"
	"net/http"

	keyvaluestore "github.com/defoeam/kvs/kv"
	"github.com/gin-gonic/gin"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func newKeyValue(key string, value string) *KeyValue {
	kv := KeyValue{
		Key:   key,
		Value: value,
	}

	return &kv
}

/*
handleGetAll returns a gin.HandlerFunc that retrieves all key-values in the store.

Example request:

	GET /items

Example response:

	[
		{
			"key": "age",
			"value": "23"
		},
		{
			"key": "name",
			"value": "Tom"
		}
	]
*/
func handleGetAll(kv *keyvaluestore.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get values for the store.
		data := kv.GetAll()

		// Build the response.
		var res []KeyValue
		res = make([]KeyValue, 0)

		for k, v := range data {
			res = append(res, *newKeyValue(k, v))
		}

		// Serialize response
		ctx.JSON(http.StatusOK, res)
	}
}

/*
handleGet returns a gin.HandlerFunc that retrieves the value associated with the provided key from the KeyValueStore.
If the key is missing or not found, it responds with an appropriate HTTP status and error message.

Example Request:

	GET /items/:key

Example Response:

	{
	  "key": "exampleKey",
	  "value": "exampleValue"
	}
*/
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
		res := newKeyValue(key, val)

		// Serialize response
		ctx.JSON(http.StatusAccepted, res)
	}
}

/*
handleSet returns a gin.HandlerFunc that sets a key-value pair in the KeyValueStore.
The key and value are provided in the JSON request body. It responds with the created key-value pair.

Example Request:

	POST /items
	{
	   "key": "exampleKey",
	   "value": "exampleValue"
	}

Example Response:

	{
	   "key": "exampleKey",
	   "value": "exampleValue"
	}
*/
func handleSet(kv *keyvaluestore.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Define key/value request structure
		var req KeyValue

		// Bind JSON to key/value request structure
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

// handleClear removes all items from the store.
func handleClear(kv *keyvaluestore.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		kv.Clear()
		ctx.Status(http.StatusNoContent)
	}
}

/*
handleGetKeys gets all keys from the store.

Example Request:

	GET /keys

Example Response:

	[
		"exampleKey1",
		"exampleKey2"
	]
*/

func handleGetKeys(kv *keyvaluestore.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keys := kv.GetKeys()
		ctx.JSON(http.StatusOK, keys)
	}
}

/*
handleGetValues returns all values from the store.

Example Request:

	GET /values

Example Response:

	[
		"exampleValue1",
		"exampleValue2"
	]
*/
func handleGetValues(kv *keyvaluestore.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		values := kv.GetValues()
		ctx.JSON(http.StatusOK, values)
	}
}

func main() {
	// Create a new instance of KeyValueStore.
	kv := keyvaluestore.NewKeyValueStore()

	// Setup gin engine
	router := gin.Default()

	// GET endpoints
	router.GET("/items", handleGetAll(kv))
	router.GET("/items/:key", handleGet(kv))
	router.GET("/keys", handleGetKeys(kv))
	router.GET("/values", handleGetValues(kv))

	// POST endpoints
	router.POST("/items", handleSet(kv))

	// DELETE endpoints
	router.DELETE("/items", handleClear(kv))

	// Get address
	port := 8080
	addr := fmt.Sprintf(":%d", port)

	// Start server
	log.Printf("Starting key-value store on http://localhost%s\n", addr)
	log.Fatal(router.Run(addr))
}
