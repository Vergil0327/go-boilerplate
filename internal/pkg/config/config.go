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
	PrintConfig bool
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
