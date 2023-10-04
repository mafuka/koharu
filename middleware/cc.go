// middleware 包提供了一组中间件，用于处理客户端上报的事件。是机器人功能的核心实现。
package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/lib/event"
	"github.com/kwaain/nakisama/lib/logger"
	"go.uber.org/zap"
)

// CustomContext 是一个自定义的上下文类型。
// 它含了所有需要在中间件之间共享的事件，用于在中间件之间以类型安全的方式共享这些事件。
//
// C: 23-10-02 22:46;
// U: 23-10-03 00:18.
type CustomContext struct {
	*gin.Context
	MessageEvent *event.MessageEvent
	RequestEvent *event.RequestEvent
	NoticeEvent  *event.NoticeEvent
	MetaEvent    *event.MetaEvent
}

// NewCustomContext 创建一个新的 CustomContext。
// 它接受一个 *gin.Context，并将其包装在 CustomContext 中。
//
// C: 23-10-02 22:47;
// U: 23-10-03 00:18.
func NewCustomContext(c *gin.Context) *CustomContext {
	return &CustomContext{
		Context: c,
	}
}

// WithCustomContext 是一个中间件，用于初始化 CustomContext 并将其存入 gin.Context。
//
// C: 23-10-03 00:19;
// U: 23-10-03 01:23.
func WithCustomContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		cc := NewCustomContext(c)
		c.Set("CustomContextKey", cc)

		c.Next()
	}
}

// UseCustomMiddlewares 是一个中间件，从 gin.Context 中检索 CustomContext 并传递给各中间件。
// 它以闭包方式实现，以保证类型安全，
//
// C: 23-10-03 00:21;
// U: 23-10-03 00:56.
func UseCustomMiddlewares() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("CustomContextKey")
		if !exists {
			err := fmt.Errorf("CustomContext 不存在")
			logger.Error("中间件链初始化错误", zap.Error(err))
			return
		}

		cc, ok := value.(*CustomContext)
		if !ok {
			err := fmt.Errorf("CustomContext 类型转换失败")
			logger.Error("中间件链初始化错误", zap.Error(err))
			return
		}

		// 自定义中间件链
		Handler(cc)(c)
		Whitelist(cc)(c)
		OffTopic(cc)(c)
	}
}
