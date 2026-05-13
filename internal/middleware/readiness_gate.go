package middleware

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type ReadinessGate struct {
	ready atomic.Bool
}

func NewReadinessGate() *ReadinessGate {
	g := &ReadinessGate{}
	g.ready.Store(true)
	return g
}

func (g *ReadinessGate) SetReady(v bool) {
	g.ready.Store(v)
}

func (g *ReadinessGate) IsReady() bool {
	return g.ready.Load()
}

func ReadinessGateMiddleware(readinesGate *ReadinessGate) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !readinesGate.IsReady() {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "service unavailable",
			})
			return
		}

		c.Next()
	}
}
