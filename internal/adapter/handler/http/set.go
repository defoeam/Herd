package http

import (
	"net/http"

	"github.com/defoeam/kvs/internal/core/domain"
	"github.com/defoeam/kvs/internal/core/port"
	"github.com/gin-gonic/gin"
)

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
func handleSet(kv port.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Define key/value request structure
		var req domain.KeyValue

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
