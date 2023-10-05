package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/middleware/event"
	"github.com/kwaain/nakisama/middleware/offtopic"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(event.Handler())
	router.Use(offtopic.Handler())

	router.POST("/", func(c *gin.Context) { c.String(http.StatusOK, "OK") })

	return router
}
