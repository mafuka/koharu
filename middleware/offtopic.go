package middleware

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/lib/api"
	"github.com/kwaain/nakisama/lib/logger"
	"go.uber.org/zap"
)

// OffTopic 标记一条跑题消息，并警告发言用户。
// 使用 /ot 回复一条消息即可触发。仅管理员有权使用。
//
// C: 23-10-02 hh:mm;
// U: 23-10-03 15:28.
func OffTopic(cc *CustomContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 非消息事件不适用本中间件
		if cc.MessageEvent == nil {
			logger.Debug("已忽略非消息事件", zap.Any("messageEvent", cc.MessageEvent))

			c.Next()
			return
		}

		e := cc.MessageEvent

		// 非 /ot 指令不适用本中间件
		if !strings.Contains(e.RawMessage, "/ot") {
			logger.Debug("已忽略非目标命令消息", zap.Any("messageEvent", e))

			c.Next()
			return
		}

		// 指令不合法则报错

		// 1. 没有回复
		if !strings.Contains(e.RawMessage, "[CQ:reply,id=") {
			logger.Debug("指令不合法：没有回复的消息", zap.Any("rawMessage", e.RawMessage))

			c.Abort()
			return
		}

		// 2. 多余参数
		// TODO: 23-10-03 17:42 懒得做了，参数多余就忽略好了

		// TODO: 23-10-03 15:30 上面的逻辑以后可能会大量复用，应该考虑抽离成“事件类型filter”和“指令filter”两个工具。
		// TODO: 23-10-03 15:33 正好还能优化下接口能力，这样每次写中间件只需要关注功能逻辑就行。
		// TODO: 23-10-03 16:52 也可以大胆一些，改中间件链，做个“预处理”的环节。

		// 获取离题消息 ID
		m := regexp.MustCompile(`id=(-?\d+)`).FindStringSubmatch(e.RawMessage)
		if !(len(m) > 1) {
			logger.Debug("无法获取离题消息 ID：正则过滤失败", zap.Any("rawMessage", e.RawMessage), zap.Any("regexResult", m))

			c.Abort()
			return
		}
		otMsgID, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			logger.Debug("无法获取离题消息 ID：类型转换失败", zap.Error(err))

			c.Abort()
			return
		}

		// 警告消息构造
		warnStr := generateWarnStr(otMsgID)

		sentMsg := api.GroupMsg{
			GroupID: e.GroupID,
			Message: warnStr,
		}

		// 发送警告消息
		sentMsgID, err := api.SendGroupMsg(sentMsg)
		if err != nil {
			logger.Debug("群聊消息发送失败", zap.Error(err))

			c.Abort()
			return
		}

		logger.Debug("消息发送成功", zap.Float64("messageID", sentMsgID)) // debug

		c.Next()
	}
}

// generateWarnStr 生成离题警告的消息内容。msgID 是被警告的消息，userID 是其发言者。
//
// TODO: 23-10-03 16:14 以后这种格式化的文本应该不少，感觉应该做个模版系统，用 text/template 包做定义。抽空改一下。
func generateWarnStr(msgID float64) (warnStr string) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[CQ:reply,id=%f]\n", msgID)) // CQ 码 回复
	builder.WriteString("======================\n")
	builder.WriteString("💡【OFF-TOPIC WARNING】\n")
	builder.WriteString("欢迎您参与讨论！\n")
	builder.WriteString("但您的发言好像离题了。\n")
	builder.WriteString("请保持话题相关，感谢理解😘\n")
	builder.WriteString("======================\n")
	builder.WriteString(fmt.Sprintf("RequestID: %f", msgID)) // debug information

	return builder.String()
}
