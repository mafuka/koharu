package middleware

import "github.com/gin-gonic/gin"

// OffTopic 标记一条跑题消息，并警告发言用户。
// 使用 /ot 回复一条消息即可触发。仅管理员有权使用。
//
// TODO: 困死了，明天再写逻辑
//
// C: 23-10-02 hh:mm;
// U: 23-10-03 01:30.
func OffTopic() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
