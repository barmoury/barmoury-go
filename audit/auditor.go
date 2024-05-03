package audit

import (
	"reflect"
	"time"

	"github.com/barmoury/barmoury-go/cache"
	"github.com/barmoury/barmoury-go/util"
)

type auditorInterface[T any] interface {
	Flush()
	PreAudit(Audit[T])
	GetCache() cache.Cache[Audit[T]]
}

type Auditor[T any] struct {
	BufferSize      uint64
	DateLastFlushed time.Time
	auditorInterface[T]
}

func NewAuditor[T any]() Auditor[T] {
	return Auditor[T]{
		BufferSize:      0,
		DateLastFlushed: time.Now(),
	}
}

func TriggerAudit_Reflect[T any, A any](t T, a A) {
	f := util.GetDeclaredSurefireMethod(t, "Flush")
	m := util.GetDeclaredSurefireMethod(t, "PreAudit")
	c := util.GetDeclaredSurefireMethod(t, "GetCache")
	m.Call([]reflect.Value{reflect.ValueOf(a)})
	bufferSize := util.GetDeclaredFieldValueAsUint(t, "BufferSize") + 1
	dateLastFlushed := util.GetDeclaredFieldValueAsTime(t, "DateLastFlushed")
	if dateLastFlushed.IsZero() {
		dateLastFlushed = time.Now()
		util.SetFieldValue(t, "DateLastFlushed", dateLastFlushed)
	}
	util.SetFieldValue(t, "BufferSize", bufferSize)
	cached := c.Call(nil)[0]
	if util.CacheWriteAlong(bufferSize, dateLastFlushed, cached.Interface().(cache.Cache[A]), a) {
		util.SetFieldValue(t, "BufferSize", uint64(0))
		util.SetFieldValue(t, "DateLastFlushed", time.Now())
		f.Call(nil)
	}
}

func (ad *Auditor[T]) Audit(t auditorInterface[T], a Audit[T]) {
	t.PreAudit(a)
	ad.BufferSize++
	if util.CacheWriteAlong(ad.BufferSize, ad.DateLastFlushed, t.GetCache(), a) {
		ad.BufferSize = 0
		ad.DateLastFlushed = time.Now()
		t.Flush()
	}
}
