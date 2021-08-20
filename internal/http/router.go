package http

import "github.com/gin-gonic/gin"

func NewRouter(handler *gin.Engine) {
	// Routes.
	h := handler.Group("/api")
	{
		newDumpRoutes(h)
	}
}

