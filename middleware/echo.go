package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/api"
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
func Echo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var message Message
		err := c.ShouldBindJSON(&message)
		if err != nil {
			logger.Error("解析请求体失败", zap.Error(err))
			return
		}

		if message.MessageType == "private" && strings.HasPrefix(message.Message, "/echo") {
			userID := message.UserID
			content := message.Message
			msgID, err := api.SendPrivateMsg(userID, 0, content, false)
			if err != nil {
				logger.Error("发送私聊消息失败", zap.Error(err))
				return
			}

			logger.Info("发送私聊消息", zap.Float64("messageID", msgID))
		}
	}
}
