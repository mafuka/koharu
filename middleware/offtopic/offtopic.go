// offtopic 标记一条跑题消息，并警告发言用户。
// 使用 /ot 回复一条消息即可触发。仅管理员有权使用。
package offtopic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/lib/api"
	"github.com/kwaain/nakisama/lib/logger"
	"github.com/kwaain/nakisama/middleware/event"
	"go.uber.org/zap"
)

// Handler 进行事件预处理
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rID := c.MustGet("requestID").(string)

		logger.Debug("进入中间件", zap.String("requestID", rID))

		e, ok := c.MustGet("event").(event.MessageEvent)
		if !ok {
			logger.Warn("拒绝处理该事件: 非消息事件", zap.String("requestID", rID), zap.Any("event", e))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Next()
			return
		}

		if e.MessageType != "group" {
			logger.Warn("拒绝处理该事件: 非群组会话", zap.String("requestID", rID), zap.Any("messageType", e.MessageType))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Next()
			return
		}

		if !strings.Contains(e.RawMessage, "/ot") {
			logger.Warn("拒绝处理该事件: 命令不匹配", zap.String("requestID", rID), zap.Any("rawMessage", e.RawMessage))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Next()
			return
		}

		if !strings.Contains(e.RawMessage, "[CQ:reply,id=") {
			logger.Warn("命令不合法: 缺少回复", zap.String("requestID", rID), zap.Any("rawMessage", e.RawMessage))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Next()
			return
		}

		// 获取 ot 消息 ID
		m := regexp.MustCompile(`id=(-?\d+)`).FindStringSubmatch(e.RawMessage)
		if !(len(m) > 1) {
			logger.Error("无法获取 ot 消息 ID: 正则过滤失败", zap.String("requestID", rID), zap.String("rawMessage", e.RawMessage), zap.Any("regexResult", m))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Abort()
			return
		}
		otMsgID, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			logger.Debug("无法获取 ot 消息 ID: 类型转换失败", zap.Error(err))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Abort()
			return
		}

		// 警告消息构造
		warnStr := generateWarnStr(rID, otMsgID)

		// 发送警告消息
		sentMsg := api.GroupMsg{
			GroupID: e.GroupID,
			Message: warnStr,
		}

		sentMsgID, err := api.SendGroupMsg(sentMsg)
		if err != nil {
			logger.Error("群聊消息发送失败", zap.String("requestID", rID), zap.Error(err))

			logger.Debug("离开中间件", zap.String("requestID", rID))
			c.Abort()
			return
		}

		logger.Info("消息发送成功", zap.String("requestID", rID), zap.Float64("messageID", sentMsgID))

		logger.Debug("离开中间件", zap.String("requestID", rID))
		c.Abort()
	}
}

// generateWarnStr 生成离题警告的消息内容。msgID 是被警告的消息，userID 是其发言者。
//
// TODO: 23-10-03 16:14 以后这种格式化的文本应该不少，感觉应该做个模版系统，用 text/template 包做定义。抽空改一下。
func generateWarnStr(requestID string, msgID float64) (warnStr string) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[CQ:reply,id=%f]\n", msgID)) // CQ 码 回复
	builder.WriteString("======================\n")
	builder.WriteString("💡【OFF-TOPIC WARNING】\n")
	builder.WriteString("欢迎您参与讨论！\n")
	builder.WriteString("但您的发言好像离题了。\n")
	builder.WriteString("请保持话题相关，感谢理解😘\n")
	builder.WriteString("======================\n")
	builder.WriteString(fmt.Sprintf("RequestID: %s", requestID)) // debug information

	return builder.String()
}
