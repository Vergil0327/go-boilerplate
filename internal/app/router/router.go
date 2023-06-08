package router

import (
	"boilerplate/internal/app/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var _ IRouter = (*Router)(nil)

var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

type Router struct {
}

func (r *Router) Register(app *gin.Engine) error {
	g := app.Group("/api")

	g.Use(middlewares.RateLimiterMiddleware())

	v1 := g.Group("/v1")
	{
		pub := v1.Group("/pub")
		pub.GET("/helloworld", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "hello world"})
		})
	}

	return nil
}

func (r *Router) Prefixes() []string {
	return []string{"/api/"}
}
