package middleware

import "github.com/gin-gonic/gin"

// Whoami 只会复述用户 QQ 号
func Whoami() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
