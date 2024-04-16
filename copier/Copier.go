package copier

import (
	"reflect"
	"strings"

	"github.com/barmoury/barmoury-go/util"
)

func Copy(target interface{}, sources ...interface{}) interface{} {
	t := util.GetFieldNonPtrType(target)
	for i := 0; i < t.NumField(); i++ {
		copyField(t.Field(i), target, sources...)
	}
	return target
}

func copyField(field reflect.StructField, target interface{}, sources ...interface{}) {
	cpt := field.Tag.Get("copy_property")
	if strings.Contains(cpt, "ignore") {
		return
	}
	fieldName := field.Name
	value, ok := findUsableValue(field, sources...)
	if !ok {
		return
	}
	if value.IsZero() && !strings.Contains(cpt, "use_zero_value") {
		return
	}
	util.SetFieldValue(target, fieldName, value)
}

func findUsableValue(field reflect.StructField, sources ...interface{}) (reflect.Value, bool) {
	ok := false
	name := field.Name
	var lv reflect.Value
	var value reflect.Value
	for _, source := range sources {
		if source == nil {
			continue
		}
		v, okk := util.GetDeclaredFieldValue(source, name)
		if !okk {
			continue
		}
		if v.Type() != field.Type {
			continue
		}
		lv = v
		ok = true
		if !v.IsZero() {
			break
		}
	}
	value = lv
	return value, ok
}
