package http

import (
	"net/http"

	"github.com/defoeam/kvs/internal/core/domain"
	"github.com/defoeam/kvs/internal/core/port"
	"github.com/gin-gonic/gin"
)

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
func handleGet(kv port.KeyValueStore) gin.HandlerFunc {
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
		res := domain.NewKeyValue(key, val)

		// Serialize response
		ctx.JSON(http.StatusAccepted, res)
	}
}
