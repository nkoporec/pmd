package http

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine) {
	// Routes.
	h := handler.Group("/api")
	{
		newApiRoutes(h)
	}
}

func newApiRoutes(handler *gin.RouterGroup) {
	handler.POST("/dump", dump)
}

