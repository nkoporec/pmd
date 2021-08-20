package http

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	tson "github.com/skanehira/tson/lib"
)

func dump(c *gin.Context) {
	request,err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	go tson.Edit(request)
	return
}
