package app

import (
	"boilerplate/internal/pkg/config"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type options struct {
	ConfigFile string
}
type Option func(*options)

func SetConfigFile(s string) Option {
	return func(opt *options) {
		opt.ConfigFile = s
	}
}

func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, setter := range opts {
		setter(&o)
	}

	config.MustLoad(o.ConfigFile)
	config.PrintWithJSON()

	clean := func() {}
	return clean, nil
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
