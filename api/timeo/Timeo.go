package timeo

import (
	"barmoury/util"
	"time"
)

func Resolve(model interface{}) {
	modelId, ok := util.GetDeclaredFieldValue(model, "Id")
	if !ok {
		return
	}
	if modelId.Uint() <= 0 {
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
