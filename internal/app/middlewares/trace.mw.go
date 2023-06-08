package middlewares

import (
	"boilerplate/internal/pkg/logger"
	"boilerplate/internal/pkg/trace"
	"context"

	"github.com/gin-gonic/gin"
)

type (
	traceIDKey struct{}
)

func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func FromTraceIDContext(ctx context.Context) (string, bool) {
	if v := ctx.Value(traceIDKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s, s != ""
		}
	}
	return "", false
}

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = trace.NewTraceID()
		}

		ctx := NewTraceIDContext(c.Request.Context(), traceID)
		ctx = logger.NewTraceIDContext(ctx, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Trace-Id", traceID)

		c.Next()
	}
}
