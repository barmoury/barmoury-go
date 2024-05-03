package meta

import "github.com/barmoury/barmoury-go/util"

func FindAnnotation(as map[string]any, n string) (any, bool) {
	v, ok := as[n]
	if ok {
		return v, true
	}
	return nil, false
}

func GetAnnotation[T any](as map[string]any, n string) T {
	v, ok := as[n]
	if ok {
		return v.(T)
	}
	return v.(T)
}

func FindAnnotationInMap(as []map[string]any, n string) (any, bool) {
	for _, d := range as {
		if d["name"] == n {
			return d, true
		}
	}
	return nil, true
}

func processAnnotationsFromMethod(i interface{}, name string, cb func(string, any)) {
	if am := util.GetDeclaredMethodValue(i, name); am.IsValid() {
		util.TranverseIterable(am.Call(nil)[0].Interface(), func(k any, v any) {
			cb(k.(string), v)
		})
	}
}

func GetAnnotationsFromMethod(i interface{}, name string) map[string]any {
	annos := make(map[string]any)
	processAnnotationsFromMethod(i, name, func(s string, a any) {
		annos[s] = a
	})
	return annos
}

func GetAnnotationsFromMethods(i interface{}, names []string) map[string]any {
	annos := make(map[string]any)
	for _, name := range names {
		processAnnotationsFromMethod(i, name, func(s string, a any) {
			annos[s] = a
		})
	}
	return annos
}

func GetAnnotationsFromAttributesAnnotationsMethods(i interface{}) map[string]any {
	return GetAnnotationsFromMethods(i, []string{"DefaultAttributesAnnotations", "AttributesAnnotations"})
}

func GetAttributesAnnotations[T any](annos map[string]any, method string, name string) (T, bool) {
	var n T
	mma, ok := annos[method].(map[string]any)
	if !ok {
		return n, false
	}
	if ea_, ok := FindAnnotation(mma, name); ok {
		return ea_.(T), true
	}
	return n, false
}
