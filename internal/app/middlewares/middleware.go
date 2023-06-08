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
