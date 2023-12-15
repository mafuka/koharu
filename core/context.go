package core

import "github.com/gin-gonic/gin"

type Context struct {
	gContext *gin.Context
	Event
}

func (c *Context) GetRawData() ([]byte, error) {
	return c.gContext.GetRawData()
}

func (c *Context) EventIs(v EventType) bool {
	return c.Event.GetType() == v
}

func (c *Context) Next() {
	c.gContext.Next()
}

func (c *Context) Abort() {
	c.gContext.Abort()
}
