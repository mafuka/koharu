package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/lib/event"
	"github.com/kwaain/nakisama/lib/logger"
	"go.uber.org/zap"
)

// Handler 处理客户端的上报事件并将其存储在 CustomContext 中。
//
// C: 23-10-01 hh:mm;
// U: 23-10-03 01:28.
func Handler(cc *CustomContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取并解析原始数据
		r, err := c.GetRawData()
		if err != nil {
			logger.Error("获取新上报事件的原始数据失败", zap.Error(err))
			return
		}

		var newEvent event.EventGeneral
		err = json.Unmarshal(r, &newEvent)
		if err != nil {
			logger.Error("解析新上报事件失败", zap.Error(err))
			return
		}

		// 根据事件类型，解析并存储事件
		switch newEvent.PostType {
		case "message":
			var messageEvent event.MessageEvent
			err := json.Unmarshal(r, &messageEvent)
			if err != nil {
				logger.Error("解析消息事件失败", zap.Error(err))
				return
			}

			cc.MessageEvent = &messageEvent

		case "message_sent":
			logger.Warn("暂不支持处理机器人自己发送的消息")
			return

		case "request":
			var requestEvent event.RequestEvent
			err := json.Unmarshal(r, &requestEvent)
			if err != nil {
				logger.Error("解析请求事件失败", zap.Error(err))
				return
			}

			cc.RequestEvent = &requestEvent

		case "notice":
			var noticeEvent event.NoticeEvent
			err := json.Unmarshal(r, &noticeEvent)
			if err != nil {
				logger.Error("解析通知事件失败", zap.Error(err))
				return
			}

			cc.NoticeEvent = &noticeEvent

		case "meta_event":
			var metaEvent event.MetaEvent
			err := json.Unmarshal(r, &metaEvent)
			if err != nil {
				logger.Warn("暂不支持处理元数据事件")
				return
			}

			cc.MetaEvent = &metaEvent
		default:
			err := fmt.Errorf("未知的事件类型 %s", newEvent.PostType)
			logger.Error("解析新上报事件失败", zap.Error(err))
		}

		// debug
		if cc.MessageEvent != nil {
			logger.Debug("消息事件解析成功", zap.Any("messageEvent", cc.MessageEvent))
		}
		if cc.RequestEvent != nil {
			logger.Debug("请求事件解析成功", zap.Any("requestEvent", cc.RequestEvent))
		}
		if cc.NoticeEvent != nil {
			logger.Debug("通知事件解析成功", zap.Any("noticeEvent", cc.NoticeEvent))
		}
		if cc.MetaEvent != nil {
			logger.Debug("元数据事件解析成功", zap.Any("metaEvent", cc.MetaEvent))
		}

		c.Next()
	}
}
