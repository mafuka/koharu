package middleware

import "github.com/gin-gonic/gin"

// Whoami 只会复述用户 ID
// TODO: 也是先开个头，明天再写逻辑
//
// C 23-09-dd hh-mm;
// U 23-10-03 01:32.
func Whoami() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
