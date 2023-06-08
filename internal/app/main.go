package app

import (
	"boilerplate/internal/pkg/config"
	"boilerplate/internal/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type options struct {
	ConfigFile string
	Version    string
}
type Option func(*options)

func SetConfigFile(s string) Option {
	return func(opt *options) {
		opt.ConfigFile = s
	}
}

func SetVersion(s string) Option {
	return func(opt *options) {
		opt.Version = s
	}
}

func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, setter := range opts {
		setter(&o)
	}

	config.MustLoad(o.ConfigFile)
	config.PrintWithJSON()

	logger.WithContext(ctx).Printf("Start server,#version %s,#pid %d", o.Version, os.Getpid())

	cleanLogger, err := InitLogger()
	if err != nil {
		return nil, err
	}

	injector, err := BuildInjector()
	if err != nil {
		return nil, err
	}

	cleanHTTP := InitHTTP(ctx, injector.Engine)

	return func() {
		cleanLogger()
		cleanHTTP()
	}, nil
}

// SIGHUP   Hangup detected on controlling terminal or death of controlling process
// SIGINT   Interrupt from keyboard. equals os.Interrupt
// SIGQUIT  Quit from keyboard
// SIGTERM  Termination signal
func Run(ctx context.Context, opts ...Option) error {
	state := 1

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)

	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return nil
	}

EXIT:
	for {
		sig := <-sc

		switch sig {
		case syscall.SIGHUP:
		case os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM:
			state = 0
			break EXIT
		default:
			break EXIT
		}
	}

	cleanFunc()

	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
