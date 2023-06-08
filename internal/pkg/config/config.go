package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/koding/multiconfig"
)

var (
	C    = new(Config)
	once sync.Once
)

type Config struct {
	RunMode     string
	PrintConfig bool

	HTTP HTTP
	CORS CORS
	GZIP GZIP
}

type HTTP struct {
	Host               string
	Port               int
	CertFile           string
	KeyFile            string
	ShutdownTimeout    int
	MaxContentLength   int64
	MaxReqLoggerLength int `default:"1024"`
	MaxResLoggerLength int
	ReadTimeout        int
	WriteTimeout       int
	IdleTimeout        int
}

type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type GZIP struct {
	Enable             bool
	ExcludedExtentions []string
	ExcludedPaths      []string
}

func MustLoad(filePaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, fpath := range filePaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}

		mconfig := &multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}

		mconfig.MustLoad(C)
	})
}

func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", "\t")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}

		os.Stdout.WriteString(fmt.Sprintln(string(b)))
	}
}
