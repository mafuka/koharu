package core

import "github.com/gin-gonic/gin"

type Context struct {
	Event
	*Client
	context *gin.Context
}

func (c *Context) GetRawData() ([]byte, error) {
	return c.context.GetRawData()
}

func (c *Context) Next() {
	c.context.Next()
}

func (c *Context) Abort() {
	c.context.Abort()
}
