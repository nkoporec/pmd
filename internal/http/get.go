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
		c.JSON(http.StatusOK, gin.H{"data": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
