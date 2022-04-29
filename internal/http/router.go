package http

import (
	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
)

func Cache(messages chan interface{}) gin.HandlerFunc {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		c.Set("cache", cache)
		c.Set("messages", messages)
		c.Next()
	}
}

func NewRouter(handler *gin.Engine, messages chan interface{}) {
	handler.Use(Cache(messages))

	// Routes.
	h := handler.Group("/api")
	{
		newApiRoutes(h)
	}
}

func newApiRoutes(handler *gin.RouterGroup) {
	handler.POST("/dump", dump)
}
