package http

import (
	"net/http"

	"github.com/defoeam/kvs/internal/core/port"
	"github.com/gin-gonic/gin"
)

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

func handleGetKeys(kv port.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keys := kv.GetKeys()
		ctx.JSON(http.StatusOK, keys)
	}
}
