package tracer

import (
	"github.com/gin-gonic/gin"
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

func MustFromContext(c *gin.Context) *Trace {
	t := c.MustGet("trace")
	return t.(*Trace)
}

// ToContext adds the Span to the gin.Context.
//
//	trace.EndSpan().ToContext(c)
func (s *Span) ToContext(c *gin.Context) *Span {
	c.Set("span", s)
	return s
}
