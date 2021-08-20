package http

import "github.com/gin-gonic/gin"

func newDumpRoutes(handler *gin.RouterGroup) {
	handler.POST("/dump", dump)
}

func dump(c *gin.Context) {
	c.JSON(200, gin.H {
		"message" : "dump",
	})
}
