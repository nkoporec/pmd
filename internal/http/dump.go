package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nkoporec/dump/internal/rpc"
)

func dump(c *gin.Context) {
	rpc.SocketHandler(c.Writer, c.Request)
}
