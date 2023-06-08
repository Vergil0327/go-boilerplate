package middlewares

import (
	"boilerplate/internal/pkg/config"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func GzipMiddleware() gin.HandlerFunc {
	cfg := config.C.GZIP
	return gzip.Gzip(gzip.BestCompression,
		gzip.WithExcludedExtensions(cfg.ExcludedExtentions),
		gzip.WithExcludedPaths(cfg.ExcludedPaths),
	)
}
