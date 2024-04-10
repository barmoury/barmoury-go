package copier

import (
	"barmoury/util"
	"reflect"
)

func Copy(target interface{}, sources ...interface{}) interface{} {
	t := util.GetFieldNonPtrType(target)
	for i := 0; i < t.NumField(); i++ {
		copyField(t.Field(i), target, sources...)
	}
	return target
}

func copyField(field reflect.StructField, target interface{}, sources ...interface{}) {
	fieldName := field.Name
	value, ok := findUsableValue(fieldName, sources...)
	if !ok {
		return
	}
	util.SetFieldValue(target, fieldName, value)
}

func findUsableValue(name string, sources ...interface{}) (reflect.Value, bool) {
	ok := false
	var value reflect.Value
	for _, source := range sources {
		if source == nil {
			continue
		}
		v, okk := util.GetDeclaredFieldValue(source, name)
		if !okk {
			continue
		}
		ok = true
		value = v
		break
	}
	return value, ok
}
