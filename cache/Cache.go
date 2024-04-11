package cache

type Cache[T any] interface {
	Cache(T)
	GetCached() []T
	MaxBufferSize() uint64
	IntervalBeforeFlush() uint64
}
