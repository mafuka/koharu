// middleware 包含了一些中间件函数，用于处理消息。
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/lib/api"
	"github.com/kwaain/nakisama/lib/logger"
	"go.uber.org/zap"
)

type Message struct {
	Time        int64  `json:"time"`
	SelfID      int64  `json:"self_id"`
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	MessageID   int32  `json:"message_id"`
	UserID      int64  `json:"user_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int    `json:"font"`
}

// Echo 是一个中间件函数，会复述私聊消息。
//
// Deprecated: 这是刚开始用中间件做消息处理时测试用的，现在已经没用了
//
// C: 23-09-04 19:mm;
// U: 23-10-03 14:41.
func Echo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var message Message
		err := c.ShouldBindJSON(&message)
		if err != nil {
			logger.Error("解析请求体失败", zap.Error(err))
			return
		}

		if message.MessageType == "private" && strings.HasPrefix(message.Message, "/echo") {
			msg := api.PrivateMsg{
				UserID:  message.UserID,
				Message: message.Message,
			}

			msgID, err := api.SendPrivateMsg(msg)
			if err != nil {
				logger.Error("发送私聊消息失败", zap.Error(err))
				return
			}

			logger.Info("发送私聊消息", zap.Float64("messageID", msgID))
		}
	}
}
