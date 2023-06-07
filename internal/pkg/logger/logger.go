package logger

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Define logrus alias
type (
	Logger = logrus.Logger
	Entry  = logrus.Entry
	Hook   = logrus.Hook
	Level  = logrus.Level
)

var (
	Tracef          = logrus.Tracef
	Debugf          = logrus.Debugf
	Infof           = logrus.Infof
	Warnf           = logrus.Warnf
	Errorf          = logrus.Errorf
	Fatalf          = logrus.Fatalf
	Panicf          = logrus.Panicf
	Printf          = logrus.Printf
	SetOutput       = logrus.SetOutput
	SetReportCaller = logrus.SetReportCaller
	StandardLogger  = logrus.StandardLogger
	ParseLevel      = logrus.ParseLevel
)

// Define logger level
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// Set logger level
func SetLevel(level Level) {
	logrus.SetLevel(level)
}

type Format string

const (
	JsonFormat Format = "json"
	TextFormat Format = "text"
)

// Set logger formatter (json/text)
func SetFormatter(format Format) {
	switch format {
	case JsonFormat:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case TextFormat:
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}
}

func AddHook(hook Hook) {
	logrus.AddHook(hook)
}

const (
	ErrStackKey = "stack"
)

type (
	errStackKey struct{}
)

func NewStackContext(ctx context.Context, stack error) context.Context {
	return context.WithValue(ctx, errStackKey{}, stack)
}

func FromStackContext(ctx context.Context) error {
	if v := ctx.Value(errStackKey{}); v != nil {
		if stk, ok := v.(error); ok {
			return stk
		}
	}
	return nil
}

// Create entry from context
func WithContext(ctx context.Context) *Entry {
	fields := logrus.Fields{}

	if v := FromStackContext(ctx); v != nil {
		fields[ErrStackKey] = fmt.Sprintf("%+v", v)
	}

	return logrus.WithContext(ctx).WithFields(fields)
}
