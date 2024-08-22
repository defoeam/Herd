package http

import (
	"net/http"

	"github.com/defoeam/kvs/internal/core/domain"
	"github.com/defoeam/kvs/internal/core/port"
	"github.com/gin-gonic/gin"
)

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
func handleGetAll(kv port.KeyValueStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get values for the store.
		data := kv.GetAll()

		// Build the response.
		var res []domain.KeyValue
		res = make([]domain.KeyValue, 0)

		for k, v := range data {
			res = append(res, *domain.NewKeyValue(k, v))
		}

		// Serialize response
		ctx.JSON(http.StatusOK, res)
	}
}
