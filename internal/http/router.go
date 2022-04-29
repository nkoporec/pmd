package http

import (
	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
)

func addContext(messages chan interface{}, cch *ristretto.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cache", cch)
		c.Set("messages", messages)
		c.Next()
	}
}

func NewRouter(handler *gin.Engine, messages chan interface{}, cch *ristretto.Cache) {
	handler.Use(addContext(messages, cch))

	// Routes.
	h := handler.Group("/api")
	{
		newApiRoutes(h)
	}
}

func newApiRoutes(handler *gin.RouterGroup) {
	handler.POST("/dump", dump)
}
