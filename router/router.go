package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/middleware/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) { c.String(http.StatusOK, "/") })
	r.POST("/hook",
		handler.Handle(),
		func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) },
	)

	return r
}
