package http

import (
	"github.com/defoeam/kvs/internal/core/port"
	"github.com/gin-gonic/gin"
)

func NewHttp(kv port.KeyValueStore) *gin.Engine {
	// Setup gin engine
	app := gin.Default()

	// GET endpoints
	app.GET("/items", handleGetAll(kv))
	app.GET("/items/:key", handleGet(kv))
	app.GET("/keys", handleGetKeys(kv))
	app.GET("/values", handleGetValues(kv))

	// POST endpoints
	app.POST("/items", handleSet(kv))

	// DELETE endpoints
	app.DELETE("/items", handleClearAll(kv))
	app.DELETE("/items/:key", handleClear(kv))

	return app
}
