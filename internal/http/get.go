package http

import (
	"net/http"

	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
)

func get(c *gin.Context) {
	cache := c.MustGet("cache").(*ristretto.Cache)

	data, found := cache.Get("ddata")
	if !found {
		c.AbortWithStatusJSON(400, gin.H {
			"error": "No updates",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
