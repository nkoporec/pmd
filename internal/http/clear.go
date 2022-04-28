package http

import (
	"net/http"

	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
)

func clear(c *gin.Context) {
	var data []*RequestData
	cache := c.MustGet("cache").(*ristretto.Cache)
	cache.Set("ddata", data, 1)
	c.JSON(http.StatusOK, gin.H{"data": data})
}
