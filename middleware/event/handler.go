package event

import (
	"encoding/json"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/lib/conf"
	"github.com/kwaain/nakisama/lib/logger"
	"github.com/lithammer/shortuuid/v3"
	"go.uber.org/zap"
)

// Handle 受理客户端的上报事件并将其存储在 CustomContext 中。
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rID := shortuuid.New()
		c.Set("requestID", rID)

		logger.Debug("进入中间件", zap.String("requestID", rID))

		// 解析新事件
		r, err := c.GetRawData()
		if err != nil {
			logger.Error("无法受理新事件: 获取原始数据失败", zap.String("requestID", rID), zap.Error(err))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Abort()
			return
		}

		var e EventGeneral
		err = json.Unmarshal(r, &e)
		if err != nil {
			logger.Error("无法受理新事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Abort()
			return
		}

		// 事件归类
		switch e.PostType {
		// 消息事件
		case "message":
			var me MessageEvent
			err := json.Unmarshal(r, &me)
			if err != nil {
				logger.Error("无法受理新的消息事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

				logger.Debug("离开中间件", zap.String("requestID", rID))
				c.Abort()
				return
			}

			// 私聊白名单
			if (me.MessageType == "private") && (!slices.Contains(conf.Get().Whitelist.Private, me.UserID)) {
				logger.Info("拒绝受理新的私聊消息事件: 用户不在白名单内", zap.String("requestID", rID), zap.Int64("userID", me.UserID))

				logger.Debug("离开中间件", zap.String("requestID", rID))
				c.Abort()
				return
			}

			// 群聊白名单
			if (me.MessageType == "group") && (!slices.Contains(conf.Get().Whitelist.Group, me.GroupID)) {
				logger.Info("拒绝受理新的群聊消息事件: 群组不在白名单内", zap.String("requestID", rID), zap.Int64("groupID", me.GroupID))

				logger.Debug("离开中间件", zap.String("requestID", rID))
				c.Abort()
				return
			}

			c.Set("event", me)

			logger.Info("已受理新的消息事件", zap.String("requestID", rID), zap.Any("messageEvent", me))

		// 送出事件
		case "message_sent":
			logger.Warn("无法受理新的送出事件: 暂不支持处理送出的消息", zap.String("requestID", rID))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Abort()
			return

		// 请求事件
		case "request":
			var re RequestEvent
			err := json.Unmarshal(r, &re)
			if err != nil {
				logger.Error("无法受理新的请求事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

				logger.Debug("离开中间件", zap.String("requestID", rID))
				c.Abort()
				return
			}

			c.Set("event", re)

			logger.Info("已受理新的请求事件", zap.String("requestID", rID), zap.Any("requestEvent", re))

		// 通知事件
		case "notice":
			var ne NoticeEvent
			err := json.Unmarshal(r, &ne)
			if err != nil {
				logger.Error("无法受理新的通知事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

				logger.Debug("离开中间件", zap.String("requestID", rID))
				c.Abort()
				return
			}

			c.Set("event", ne)

			logger.Info("已受理新的通知事件", zap.String("requestID", rID), zap.Any("noticeEvent", e))

		// 元数据事件
		case "meta_event":
			var me MetaEvent
			err := json.Unmarshal(r, &me)
			if err != nil {
				logger.Error("无法受理新的元数据事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

				logger.Debug("离开中间件", zap.String("requestID", rID))
				c.Abort()
				return
			}

			c.Set("event", me)

			logger.Info("已受理新的元数据事件", zap.String("requestID", rID), zap.Any("metaEvent", me))

		// 未知事件
		default:
			logger.Warn("无法受理新事件: 未知的事件类型", zap.String("requestID", rID))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Abort()
			return
		}

		logger.Debug("离开中间件", zap.String("requestID", rID))
		c.Next()
	}
}
