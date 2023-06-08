package app

import (
	"boilerplate/internal/pkg/config"
	"boilerplate/internal/pkg/logger"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func InitLogger() (func(), error) {
	cfg := config.C.Log
	logger.SetLevel(logger.Level(cfg.Level))
	logger.SetFormatter(logger.Format(cfg.Format))

	var file *rotatelogs.RotateLogs
	if cfg.Output != "" {
		switch cfg.Output {
		case "stdout":
			logger.SetOutput(os.Stdout)
		case "stderr":
			logger.SetOutput(os.Stderr)
		case "file":
			if name := cfg.OutputFile; name != "" {
				_ = os.MkdirAll(filepath.Dir(name), 0777)

				f, err := rotatelogs.New(name+".%Y-%m-%d",
					rotatelogs.WithLinkName(name),
					rotatelogs.WithRotationTime(time.Duration(cfg.RotationTime)*time.Hour),
					rotatelogs.WithRotationCount(uint(cfg.RotationCount)))
				if err != nil {
					return nil, err
				}

				logger.SetOutput(f)
				file = f
			}
		}
	}

	return func() {
		if file != nil {
			file.Close()
		}
	}, nil
}
