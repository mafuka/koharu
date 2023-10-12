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
		// Trace
		t := tracer.NewTrace().ToContext(c)
		s := t.StartSpan()

		// Receive event
		e, err := recvEvent(c, s)
		if err != nil {
			logger.Error("Failed to receive event",
				zap.String("traceID", t.ID),
				zap.String("spanID", s.ID),
				zap.Error(err))
			t.EndSpan(s)
			c.Abort()
			return
		}
		logger.Info("New event received",
			zap.String("traceID", t.ID),
			zap.String("spanID", s.ID),
			zap.Any("event", e),
		)

		// Save

		c.Next()
	}
}

// recvEvent parses the event data and returns the event structure.
func recvEvent(c *gin.Context, s *tracer.Span) (interface{}, error) {
	r, err := c.GetRawData()
	if err != nil {
		logger.Error("Failed to get raw JSON data",

			zap.String("spanID", s.ID),
			zap.Error(err),
		)
		return nil, err
	}

	e, err := event.ParseJSON(r)
	if err != nil {
		logger.Error("Failed to parse event data",
			zap.String("spanID", s.ID),
			zap.Error(err),
		)
		return nil, err
	}

	return e, nil
}
