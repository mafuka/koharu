package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaain/nakisama/lib/event"
	"github.com/kwaain/nakisama/lib/logger"
	"github.com/kwaain/nakisama/lib/tracer"
	"go.uber.org/zap"
)

// Handle is the first in the middleware chain,
// responsible for parsing data reported via the webhook.
// Valid data will be parsed into corresponding types of structs.
func Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var k = 0

		// Trace
		t := tracer.NewTrace().ToContext(c)
		s := t.StartSpan()

		// Receive event
		e, _, err := recvEvent(c, s)
		if err != nil {
			logger.Error("Failed to receive event", zap.Error(err))
			k++
		}
		logger.Info("New event received",
			zap.String("spanID", s.ID),
			zap.Any("event", e),
		)

		// Quit
		t.EndSpan(s)
		if k > 0 {
			c.Abort()
			return
		}
		c.Next()
	}
}

// recvEvent parses the event data and returns the event structure.
func recvEvent(c *gin.Context, s *tracer.Span) (interface{}, event.Type, error) {
	r, err := c.GetRawData()
	if err != nil {
		logger.Error("Failed to get raw JSON data",
			zap.String("spanID", s.ID),
			zap.Error(err),
		)
		return nil, "", err
	}

	e, et, err := event.ParseJSON(r)
	if err != nil {
		logger.Error("Failed to parse event data",
			zap.String("spanID", s.ID),
			zap.Error(err),
		)
		return nil, "", err
	}

	return e, et, nil
}
