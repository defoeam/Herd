package http

import (
	"net/http"

	"github.com/defoeam/kvs/internal/core/domain"
	"github.com/defoeam/kvs/internal/core/port"
	"github.com/gin-gonic/gin"
)

// handleClear removes an item from the store and returns a copy of the deleted key/value pair.
func handleClear(kv port.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Param("key")

		// Check if empty key
		if key == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Key parameter is missing"})
			return
		}

		// Check if key exists and delete value
		val, ok := kv.Clear(key)
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
