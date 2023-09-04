package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.Echo())

	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "I am Nakisama, please take care of me!")
	})

	return router
}
