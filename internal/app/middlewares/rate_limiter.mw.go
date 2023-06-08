package middlewares

import (
	"boilerplate/internal/pkg/config"
	"boilerplate/internal/pkg/contextx"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"golang.org/x/time/rate"
)

func RateLimiterMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config.C

	if !cfg.RateLimiter.Enable {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": cfg.Redis.Addr,
		},
		Password: cfg.Redis.Password,
		DB:       cfg.RateLimiter.RedisDB,
	})

	limiter := redis_rate.NewLimiter(ring)
	limiter.Fallback = rate.NewLimiter(rate.Inf, 0)

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID := contextx.FromUserID(c.Request.Context())

		if userID != 0 {
			limit := cfg.RateLimiter.Count
			rate, delay, allowed := limiter.AllowMinute(fmt.Sprintf("%d", userID), limit)
			if !allowed {
				h := c.Writer.Header()
				h.Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
				h.Set("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
				delaySec := int64(delay / time.Second)
				h.Set("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))

				code := http.StatusTooManyRequests
				c.JSON(code, ErrResponse{
					Error: ErrPayload{
						Code:    code,
						Message: "too_many_requests",
					},
				})
				return
			}
		}

		c.Next()
	}
}
