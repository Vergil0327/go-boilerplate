package middlewares

import (
	"boilerplate/internal/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		path := c.Request.URL.Path

		start := time.Now()
		fields := make(map[string]interface{})
		fields["ip"] = c.ClientIP()
		fields["method"] = c.Request.Method
		fields["url"] = c.Request.URL.String()
		fields["protocol"] = c.Request.Proto
		fields["header"] = c.Request.Header
		fields["user_agent"] = c.Request.UserAgent()
		fields["content_length"] = c.Request.ContentLength

		c.Next()

		timeConsuming := time.Since(start).Nanoseconds() / 1e6
		fields["res_status"] = c.Writer.Status()
		fields["res_length"] = c.Writer.Size()

		ctx := c.Request.Context()
		entry := logger.WithContext(logger.NewTagContext(ctx, "__request__"))
		entry.
			WithFields(fields).
			Infof("[http] %s-%s-%s-%d(%dms)", path, c.Request.Method, c.ClientIP(), c.Writer.Status(), timeConsuming)
	}
}
