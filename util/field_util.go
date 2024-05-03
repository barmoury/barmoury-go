package util

import (
	"reflect"
	"time"
)

func GetFieldNonPtrType(i interface{}) reflect.Type {
	s := reflect.TypeOf(i)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	return s
}

func GetFieldPtrType(i interface{}) reflect.Type {
	s := reflect.TypeOf(i)
	return s
}

func GetFieldPtrValue(i interface{}) reflect.Value {
	s := reflect.ValueOf(i)
	if s.Kind() == reflect.Ptr {
		return s
	}
	return s.Addr()
}

func GetFieldNonPtrValue(i interface{}) reflect.Value {
	s := reflect.ValueOf(i)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	return s
}

func GetDeclaredField(i interface{}, name string) (reflect.StructField, bool) {
	t := GetFieldNonPtrType(i)
	return t.FieldByName(name)
}

func GetDeclaredFieldValue(i interface{}, name string) (reflect.Value, bool) {
	s := GetFieldNonPtrValue(i)
	v := s.FieldByName(name)
	if !v.IsValid() {
		return v, false
	}
	return v, true
}

func SetFieldValue(i interface{}, name string, value interface{}) {
	v := GetFieldNonPtrValue(i)
	f := v.FieldByName(name)
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
		if f.Type() == reflect.TypeOf(name) {
			f.SetString(value.(string))
			return
		}
	}
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	f.Set(reflect.ValueOf(value))

}

func GetDeclaredMethod(i interface{}, name string) (reflect.Method, bool) {
	t := GetFieldNonPtrType(i)
	return t.MethodByName(name)
}

func GetDeclaredMethodValue(i interface{}, name string) reflect.Value {
	t := GetFieldNonPtrValue(i)
	return t.MethodByName(name)
}

func GetPtrDeclaredMethod(i interface{}, name string) (reflect.Method, bool) {
	t := reflect.TypeOf(i)
	return t.MethodByName(name)
}

func GetDeclaredFieldValueAsUint(i interface{}, name string) uint64 {
	v := GetFieldNonPtrValue(i)
	f := v.FieldByName(name)
	value, ok := GetDeclaredFieldValue(i, name)
	if !ok {
		return uint64(0)
	}
	switch f.Kind() {
	case reflect.Uint:
	case reflect.Uint64:
		nv, ok := value.Interface().(uint64)
		if ok {
			return nv
		}
	}

	return uint64(0)
}

func GetDeclaredFieldValueAs[T any](i interface{}, name string) T {
	value, _ := GetDeclaredFieldValue(i, name)
	return value.Interface().(T)
}

func GetDeclaredFieldValueAsTime(i interface{}, name string) time.Time {
	return GetDeclaredFieldValueAs[time.Time](i, name)
}

func GetDeclaredSurefireMethod(i interface{}, name string) reflect.Value {
	t := reflect.ValueOf(i)
	return t.MethodByName(name)
}

func GetDeclaredSurefireMethodFromPtr(i interface{}, name string) reflect.Value {
	t := reflect.ValueOf(i).Elem()
	return t.MethodByName(name)
}

func GetDeclaredSurefireField(i interface{}, name string) reflect.Value {
	t := reflect.ValueOf(i)
	return t.FieldByName(name)
}

func GetDeclaredSurefireFieldFromPtr(i interface{}, name string) reflect.Value {
	t := reflect.ValueOf(i).Elem()
	return t.FieldByName(name)
}

func StructDeclaredMethod(i interface{}, name string) bool {
	t := GetFieldNonPtrType(i)
	_, ok := t.MethodByName(name)
	return ok
}

func StructImplementsDeclared(i interface{}, name string) bool {
	t := GetFieldNonPtrValue(i)
	return !t.MethodByName(name).IsZero()
}

func TranverseDeclaredMethods(i interface{}, fn func(reflect.Method, reflect.Value)) {
	t := GetFieldNonPtrType(i)
	v := GetFieldNonPtrValue(i)
	for j := 0; j < v.NumMethod(); j++ {
		fn(t.Method(j), v.Method(j))
	}
}

func GetTypeName(i interface{}) string {
	t := GetFieldNonPtrType(i)
	return t.Name()
}

func InvokeSurefireMethod(i interface{}, name string, params ...reflect.Value) []reflect.Value {
	m := GetDeclaredSurefireMethod(i, name)
	return m.Call(params)
}
