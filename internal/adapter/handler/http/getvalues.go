package http

import (
	"net/http"

	"github.com/defoeam/kvs/internal/core/port"
	"github.com/gin-gonic/gin"
)

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
func handleGetValues(kv port.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		values := kv.GetValues()
		ctx.JSON(http.StatusOK, values)
	}
}
