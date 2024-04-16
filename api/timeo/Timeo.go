package timeo

import (
	"time"

	"github.com/barmoury/barmoury-go/util"
)

func Resolve(model interface{}) {
	modelId, ok := util.GetDeclaredFieldValue(model, "Id")
	if !ok {
		return
	}
	if modelId.IsZero() {
		ResolveCreated(model)
	} else {
		ResolveUpdated(model)
	}
}

func ResolveCreated(model interface{}) {
	util.SetFieldValue(model, "CreatedAt", time.Now())
	ResolveUpdated(model)
}

func ResolveUpdated(model interface{}) {
	util.SetFieldValue(model, "UpdatedAt", time.Now())
}

func DateDiffInMinutes(a time.Time, b time.Time) uint64 {
	d := b.Sub(a)
	return uint64(d.Minutes())
}
