package tracer

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/lib/logger"
	"go.uber.org/zap"
)

/* The functions below are for integrating with other package. */

// ToContext adds the Trace to a gin.Context.
//
//	trace.NewTrace().ToContext(c)
func (t *Trace) ToContext(c *gin.Context) *Trace {
	c.Set("trace", t)
	return t
}

// Retrieves the Trace from a gin.Context.
func FromContext(c *gin.Context) (trace *Trace, exists bool) {
	t, e := c.Get("trace")
	if !e {
		return nil, false
	}
	return t.(*Trace), true
}

// ToContext adds the Span to the gin.Context.
//
//	trace.EndSpan().ToContext(c)
func (s *Span) ToContext(c *gin.Context) *Span {
	c.Set("span", s)
	return s
}

// ToLogger is a way of error handling.
// It logs the trace to Logger if err is nil, otherwise it logs it together.
// It is usually used when ending a trace:
//
//	tracer.ToLogger(trace.EndTrace)
//
// Note: Logger must be an initialized global singleton.
func ToLogger(t *Trace, err error) {
	if logger.Log == nil {
		panic("Logger has not been initialized")
	}

	if err != nil {
		logger.Error("Failed to end trace", zap.Error(err))
		return
	}

	logger.Info("Trace ended", zap.Any("trace", t))
}
