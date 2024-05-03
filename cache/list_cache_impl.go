package cache

type ListCacheImpl[T any] struct {
	cached []T
}

func NewListCacheImpl[T any]() ListCacheImpl[T] {
	return ListCacheImpl[T]{
		cached: make([]T, 0),
	}
}

func (c *ListCacheImpl[T]) MaxBufferSize() uint64 {
	return 150
}

func (c *ListCacheImpl[T]) IntervalBeforeFlush() uint64 {
	return 20
}

func (c *ListCacheImpl[T]) Cache(t T) {
	c.cached = append(c.cached, t)
}

func (c *ListCacheImpl[T]) GetCached() []T {
	resultCached := c.cached
	c.cached = make([]T, 0)
	return resultCached
}
