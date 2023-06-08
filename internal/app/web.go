package app

import (
	"boilerplate/internal/app/middlewares"
	"boilerplate/internal/app/router"
	"boilerplate/internal/pkg/config"
	"boilerplate/internal/pkg/logger"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func InitGinEngine(r router.IRouter) *gin.Engine {
	cfg := config.C
	gin.SetMode(cfg.RunMode)

	app := gin.New()

	// middlewares
	app.NoMethod(middlewares.NoMethodHandler())
	app.NoRoute(middlewares.NoRouteHandler())

	prefixes := r.Prefixes()
	app.Use(middlewares.TraceMiddleware(middlewares.NotAllowPathPrefixSkip(prefixes...)))
	app.Use(middlewares.LoggerMiddleware(middlewares.NotAllowPathPrefixSkip(prefixes...)))

	if cfg.CORS.Enable {
		app.Use(middlewares.CORSmiddleware())
	}

	if cfg.GZIP.Enable {
		app.Use(middlewares.GzipMiddleware())
	}

	r.Register(app)

	return app
}

func InitHTTP(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP
	host, port := "localhost", 3000
	addr := fmt.Sprintf("%s:%d", host, port)

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		logger.WithContext(ctx).Printf("Server is running at %s", addr)

		var err error
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			server.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.ShutdownTimeout)*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
	}
}
