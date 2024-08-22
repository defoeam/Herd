package http

import (
	"net/http"

	"github.com/defoeam/kvs/internal/core/port"
	"github.com/gin-gonic/gin"
)

// handleClearAll removes all items from the store.
func handleClearAll(kv port.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		kv.ClearAll()
		ctx.Status(http.StatusNoContent)
	}
}
