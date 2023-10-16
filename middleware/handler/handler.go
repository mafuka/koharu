package handler

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwaain/nakisama/lib/conf"
	"github.com/kwaain/nakisama/lib/event"
	"github.com/kwaain/nakisama/lib/logger"
	"go.uber.org/zap"
)

// Handle is the first in the middleware chain,
// responsible for parsing data reported via the webhook.
// Valid data will be parsed into corresponding types of structs.
func Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set("traceID", traceID)

		e, err := recvEvent(c, traceID)
		if err != nil {
			logger.ErrorT("Failed to receive event", traceID, zap.Error(err))
			c.Abort()
			return
		}
		c.Set("event", e)

		if !isWhitelisted(e, traceID) {
			logger.InfoT("Ignore non-whitelisted events", traceID)
			c.Abort()
			return
		}

		logger.InfoT("New event received", traceID, zap.Any("event", e))

		c.Next()
	}
}

// recvEvent parses the event data and returns the event structure.
func recvEvent(c *gin.Context, traceID string) (interface{}, error) {
	r, err := c.GetRawData()
	if err != nil {
		logger.ErrorT("Failed to get raw JSON data", traceID, zap.Error(err))
		return nil, err
	}

	e, err := event.ParseJSON(r)
	if err != nil {
		logger.ErrorT("Failed to parse event data", traceID, zap.Error(err))
		return nil, err
	}

	return e, nil
}

func isWhitelisted(e interface{}, traceID string) bool {
	wfriend := conf.Get().Whitelist.Friend
	wgroup := conf.Get().Whitelist.Group

	switch e := e.(type) {
	case event.FriendMsg:
		if slices.Contains(wfriend, int64(e.Sender.ID)) {
			return true
		}
	case event.GroupMsg:
		if slices.Contains(wgroup, int64(e.Sender.Group.ID)) {
			return true
		}
	default:
		return true
	}

	return false
}
