package event

import (
	"encoding/json"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwaain/nakisama/lib/conf"
	"github.com/kwaain/nakisama/lib/logger"
	"go.uber.org/zap"
)

// Handle 受理客户端的上报事件并将其存储在 CustomContext 中。
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// rID := shortuuid.New()
		rID := uuid.New().String()
		c.Set("requestID", rID)
		logger.Debug("进入中间件 event", zap.String("requestID", rID))

		// 解析新事件
		r, err := c.GetRawData()
		if err != nil {
			logger.Error("无法接收该事件: 获取原始数据失败", zap.String("requestID", rID), zap.Error(err))

			c.Abort()
			return
		}

		var e EventGeneral
		err = json.Unmarshal(r, &e)
		if err != nil {
			logger.Error("无法接收新事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

			c.Abort()
			return
		}

		logger.Info("收到新事件", zap.String("requestID", rID), zap.Any("event", e))

		// 事件归类
		switch e.PostType {
		// 消息事件
		case "message":
			var me MessageEvent
			err := json.Unmarshal(r, &me)
			if err != nil {
				logger.Error("无法受理该消息事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

				c.Abort()
				return
			}

			// 私聊白名单
			if (me.MessageType == "private") && (!slices.Contains(conf.Get().Whitelist.Private, me.UserID)) {
				logger.Info("拒绝受理该消息事件: 用户不在白名单内", zap.String("requestID", rID), zap.Int64("userID", me.UserID))

				c.Abort()
				return
			}

			// 群聊白名单
			if (me.MessageType == "group") && (!slices.Contains(conf.Get().Whitelist.Group, me.GroupID)) {
				logger.Info("拒绝受理该消息事件: 群组不在白名单内", zap.String("requestID", rID), zap.Int64("groupID", me.GroupID))

				c.Abort()
				return
			}

			c.Set("event", me)

			logger.Info("已受理该消息事件", zap.String("requestID", rID), zap.Int32("messageID", me.MessageID), zap.Int64("userID", me.UserID), zap.Int64("groupID", me.GroupID), zap.String("rawMessage", me.RawMessage))

		// 送出事件
		case "message_sent":
			logger.Warn("无法受理该送出事件: 暂不支持处理送出的消息", zap.String("requestID", rID))

			c.Abort()
			return

		// 请求事件
		case "request":
			var re RequestEvent
			err := json.Unmarshal(r, &re)
			if err != nil {
				logger.Error("无法受理该请求事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

				c.Abort()
				return
			}

			c.Set("event", re)

			logger.Info("已受理该请求事件", zap.String("requestID", rID), zap.Any("requestType", re.RequestType))

		// 通知事件
		case "notice":
			var ne NoticeEvent
			err := json.Unmarshal(r, &ne)
			if err != nil {
				logger.Error("无法受理该通知事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

				c.Abort()
				return
			}

			c.Set("event", ne)

			logger.Info("已受理该通知事件", zap.String("requestID", rID), zap.Any("noticeType", ne.NoticeType))

		// 元数据事件
		case "meta_event":
			var me MetaEvent
			err := json.Unmarshal(r, &me)
			if err != nil {
				logger.Error("无法受理该元数据事件: 解析原始数据失败", zap.String("requestID", rID), zap.Error(err))

				c.Abort()
				return
			}

			c.Set("event", me)

			logger.Info("已受理该元数据事件", zap.String("requestID", rID), zap.Any("metaEventType", me.MetaEventType))

		// 未知事件
		default:
			logger.Warn("无法受理该事件: 未知的事件类型", zap.String("requestID", rID), zap.Any("postType", e.PostType))

			c.Abort()
			return
		}

		c.Next()
	}
}
