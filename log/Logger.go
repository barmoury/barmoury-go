package log

import (
	"os"
	"time"

	"github.com/barmoury/barmoury-go/cache"
	"github.com/barmoury/barmoury-go/util"
)

type PureLogger interface {
	Log(...any)
}

type loggerInterface interface {
	Flush()
	PreLog(*Log)
	GetCache() cache.Cache[Log]
}

type Logger struct {
	BufferSize      uint64
	DateLastFlushed time.Time
	loggerInterface
}

func NewLogger() Logger {
	return Logger{
		BufferSize:      0,
		DateLastFlushed: time.Now(),
	}
}

func (ad *Logger) Log(t loggerInterface, a Log) {
	t.PreLog(&a)
	ad.BufferSize++
	if util.CacheWriteAlong(ad.BufferSize, ad.DateLastFlushed, t.GetCache(), a) {
		ad.BufferSize = 0
		ad.DateLastFlushed = time.Now()
		t.Flush()
	}
}

func formatContent(s string, args ...any) string {
	return util.StrFormat(s, args...)
}

func (ad *Logger) Verbose(t loggerInterface, s string, args ...any) {
	ad.Log(t, Log{Level: VERBOSE, Content: formatContent(s, args...)})
}

func (ad *Logger) Info(t loggerInterface, s string, args ...any) {
	ad.Log(t, Log{Level: INFO, Content: formatContent(s, args...)})
}

func (ad *Logger) Warn(t loggerInterface, s string, args ...any) {
	ad.Log(t, Log{Level: WARN, Content: formatContent(s, args...)})
}

func (ad *Logger) Trace(t loggerInterface, s string, args ...any) {
	ad.Log(t, Log{Level: TRACE, Content: formatContent(s, args...)})
}

func (ad *Logger) Error(t loggerInterface, s string, args ...any) {
	s = formatContent(s, args...) + "\n" + util.StackTraceAsString(4)
	ad.Log(t, Log{Level: ERROR, Content: s})
}

func (ad *Logger) Fatal(t loggerInterface, s string, args ...any) {
	s = formatContent(s, args...) + "\n" + util.StackTraceAsString(4)
	ad.Log(t, Log{Level: FATAL, Content: s})
	os.Exit(-1199810)
}

func (ad *Logger) Panic(t loggerInterface, s string, args ...any) {
	ad.Log(t, Log{Level: PANIC, Content: formatContent(s, args...)})
	panic(formatContent(s, args...))
}
