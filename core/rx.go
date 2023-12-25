package core

import "golang.org/x/xerrors"

// Rx initializes the Context in the Goroutine that handles the request when it receives it.
func Rx(bot *Bot) Middleware {
	return func(c *Context) {
		if err := setEvent(c); err != nil {
			Log().Error("Unable to set event for c:\n%+v", err)
			c.Abort()
		}
		setClient(bot, c)
		printEvent(c)

		c.Next()
	}
}

func setEvent(c *Context) error {
	jsonByte, err := c.GetRawData()
	if err != nil {
		return xerrors.Errorf("Unable to get raw JSON data: %w", err)
	}

	event, err := ParseEvent(jsonByte)
	if err != nil {
		return xerrors.Errorf("Failed to parse event: %w", err)
	}

	c.Event = event
	return nil
}

func setClient(b *Bot, c *Context) {
	c.Client = b.Client
}

func printEvent(c *Context) {
	switch e := c.Event.(type) {
	case *FriendMessage:
		Log().Info("[Friend] %s(%d): %v", e.Sender.Nickname, e.Sender.ID, e.MessageChain)
	case *GroupMessage:
		Log().Info("[Group] [%s(%d)] %s(%d): %+v", e.Sender.Group.Name, e.Sender.Group.ID, e.Sender.MemberName, e.Sender.ID, e.MessageChain)
	default:
		Log().Info("[%s]", e.(Event).EventType())
	}
}
