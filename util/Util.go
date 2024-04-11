package util

import (
	"fmt"
	"time"

	"github.com/barmoury/barmoury-go/cache"
)

func CacheWriteAlong[T any](bufferSize uint64, dateLastFlushed time.Time, cache cache.Cache[T], entry T) bool {
	cache.Cache(entry)
	diff := DateDiffInMinutes(dateLastFlushed, time.Now())
	return bufferSize >= cache.MaxBufferSize() || diff >= cache.IntervalBeforeFlush()
}

func DateDiffInMinutes(a time.Time, b time.Time) uint64 {
	d := b.Sub(a)
	return uint64(d.Minutes())
}

func StrFormat(str string, args ...any) string {
	return fmt.Sprintf(str, args...)
}
