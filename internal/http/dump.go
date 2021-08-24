package http

import (
	"io/ioutil"

	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
)

func dump(c *gin.Context) {
	cache := c.MustGet("cache").(*ristretto.Cache)

	request,err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	cache.Set("ddata", string(request), 1)
}
