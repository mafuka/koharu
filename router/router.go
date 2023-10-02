package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// U: 23-10-03 00:26.
	router.Use(middleware.WithCustomContext())    // 初始化自定义上下文
	router.Use(middleware.UseCustomMiddlewares()) // 使用自定义中间件

	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "I am Nakisama, please take care of me!")
	})

	return router
}
