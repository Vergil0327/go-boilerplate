package middlewares

import (
	"boilerplate/internal/pkg/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSmiddleware() gin.HandlerFunc {
	cfg := config.C.CORS
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Second * time.Duration(cfg.MaxAge),
	})
}
