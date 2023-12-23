package core

func Rx() Middleware {
	return func(c *Context) {
		r, err := c.GetRawData()
		if err != nil {
			Log().Error("Failed to get raw JSON data: %s", err)
			c.Abort()
			return
		}

		e, err := ParseEvent(r)
		if err != nil {
			Log().Error("Failed to parse event data: %s", err)
			c.Abort()
			return
		}

		c.Event = e

		switch e := e.(type) {
		case *FriendMessage:
			Log().Info("[Friend] %s(%d): %v", e.Sender.Nickname, e.Sender.ID, e.MessageChain)
		case *GroupMessage:
			Log().Info("[Group] [%s(%d)] %s(%d): %+v", e.Sender.Group.Name, e.Sender.Group.ID, e.Sender.MemberName, e.Sender.ID, e.MessageChain)
		default:
			Log().Info("[%s]", e.(Event).EventType())
		}

		c.Next()
	}
}
