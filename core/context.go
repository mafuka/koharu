package core

import "github.com/gin-gonic/gin"

type Context struct {
	Event
	*Bot
	gContext *gin.Context
}

func (c *Context) GetRawData() ([]byte, error) {
	return c.gContext.GetRawData()
}

func (c *Context) Next() {
	c.gContext.Next()
}

func (c *Context) Abort() {
	c.gContext.Abort()
}
