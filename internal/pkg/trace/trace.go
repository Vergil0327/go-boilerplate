package trace

import (
	"boilerplate/internal/pkg/logger"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var (
	nonce int64
	pid   = os.Getpid()
)

func NewTraceID() string {
	return fmt.Sprintf(
		"trace-id-%d-%s-%d",
		pid,
		time.Now().Format("2006.01.02.15.04.05.999"),
		atomic.AddInt64(&nonce, 1))
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	logger.Debugf("%v: %v\n", msg, time.Since(start))
}

// calculate execution time of function call
// usage: defer trace.Performance("Function Name or Any Label You Like") // on top of function call
func Performance(label string) {
	duration(track(label))
}
