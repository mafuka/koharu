package middleware

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/lib/logger"
	"go.uber.org/zap"
)

// Whitelist 使用白名单限制会话消息的处理。只有来自白名单会话的消息事件才会被处理。
//
// TODO: 23-10-03 01:18 有点困了，数据库部分还没写，先把白名单写死在代码里凑合用。改天补全数据库操作和管理操作。
//
// C: 23-10-02 hh:mm;
// U: 23-10-03 01:29.
func Whitelist(cc *CustomContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 本中间件仅适用消息事件
		if cc.MessageEvent == nil {
			logger.Debug("已忽略非消息事件", zap.Any("messageEvent", cc.MessageEvent))

			c.Next()
			return
		}

		e := cc.MessageEvent

		userList := []int{6025868}
		groupList := []int{345868660}

		switch e.MessageType {
		case "private":
			if !slices.Contains(userList, int(e.UserID)) {
				logger.Debug("已忽略非白名单的私聊会话", zap.Any("userId", e.UserID))

				c.Abort()
				return
			}

		case "group":
			if !slices.Contains(groupList, int(e.GroupID)) {
				logger.Debug("已忽略非白名单的群聊会话", zap.Any("groupId", e.GroupID))

				c.Abort()
				return
			}

		default:
			err := fmt.Errorf("未知的消息类型 %s", e.MessageType)
			logger.Error("白名单会话校验失败", zap.Error(err))

			c.Abort()
			return
		}
	}
}
