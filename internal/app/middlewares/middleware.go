package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	Error ErrPayload `json:"error"`
}

type ErrPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := http.StatusMethodNotAllowed
		c.JSON(code, ErrResponse{
			Error: ErrPayload{
				Code:    code,
				Message: "method_not_allowed",
			},
		})
		c.Abort()
	}
}

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := http.StatusNotFound
		c.JSON(code, ErrResponse{
			Error: ErrPayload{
				Code:    code,
				Message: "not_found",
			},
		})
	}
}

type SkipperFunc func(*gin.Context) bool

func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

func AllowPathPrefixSkip(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func NotAllowPathPrefixSkip(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}
