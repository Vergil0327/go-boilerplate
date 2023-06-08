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
	Log         Log
	RateLimiter RateLimiter
	Redis       Redis
	Gorm        Gorm
	MySQL       MySQL
	Postgres    Postgres
	Sqlite3     Sqlite3

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

type Log struct {
	Level         int
	Format        string
	Output        string
	OutputFile    string
	RotationCount int
	RotationTime  int
}

type RateLimiter struct {
	Enable  bool
	Count   int64
	RedisDB int
}

type Redis struct {
	Addr     string
	Password string
}

type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
}

type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		a.Host, a.Port, a.User, a.DBName, a.Password, a.SSLMode)
}

type Sqlite3 struct {
	Path string
}

func (a Sqlite3) DSN() string {
	return a.Path
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
