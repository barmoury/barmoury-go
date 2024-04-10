package util

import (
	"reflect"
)

func GetFieldNonPtrType(i interface{}) reflect.Type {
	s := reflect.TypeOf(i)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	return s
}

func GetFieldNonPtrValue(i interface{}) reflect.Value {
	s := reflect.ValueOf(i)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	return s
}

func GetDeclaredField(i interface{}, fieldName string) (reflect.StructField, bool) {
	t := reflect.TypeOf(i).Elem()
	return t.FieldByName(fieldName)
}

func GetDeclaredFieldValue(i interface{}, fieldName string) (reflect.Value, bool) {
	s := GetFieldNonPtrValue(i)
	v := s.FieldByName(fieldName)
	if !v.IsValid() {
		return v, false
	}
	return v, true
}

func SetFieldValue(i interface{}, fieldName string, value interface{}) {
	v := reflect.ValueOf(i).Elem()
	f := v.FieldByName(fieldName)
	if !f.CanSet() {
		return
	}
	switch value.(type) {
	case reflect.Value:
		value = value.(reflect.Value).Interface()
	}
	switch f.Kind() {
	case reflect.Int:
	case reflect.Int64:
		f.SetInt(int64(value.(int64)))
		return
	case reflect.Uint:
	case reflect.Uint64:
		f.SetUint(uint64(value.(uint64)))
		return
	case reflect.String:
		f.SetString(value.(string))
		return
	}
	f.Set(reflect.ValueOf(value))

}
