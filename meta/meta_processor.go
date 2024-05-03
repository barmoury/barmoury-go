package meta

type IAnotation interface {
	Annotations() (map[string]any, string)
}

type FieldAnotations struct {
	Classes map[string]any
}

type MethodAnotations struct {
	Parameters map[string]FieldAnotations
}

type StructAnotations struct {
	Fields           map[string]FieldAnotations
	MethodAnotations map[string]MethodAnotations
}

type InterfaceAnotations struct {
	MethodAnotations map[string]MethodAnotations
}

type PackageAnotations struct {
	Fields     map[string]FieldAnotations
	Structs    map[string]StructAnotations
	Functions  map[string]MethodAnotations
	Interfaces map[string]InterfaceAnotations
}

type Anotations struct {
	PackageAnotations map[string]PackageAnotations
}

type AnnotationProcessorOption struct {
	Logger             func(string)
	StoreInGlobalScope bool
}

func ProcessStringToAnnotationType[T any](str string) {

}

func ProcessAnnotationTypeToString[T any](t T) {

}

func AnnotationsFromSourceString(source string, options AnnotationProcessorOption) map[string]any {
	var m map[string]any
	return m
}

func ProcessAnnotationsFromFile(source string, options AnnotationProcessorOption) map[string]any {
	var m map[string]any
	return m
}

func ProcessAnnotationsFromDirectory(directory string, options AnnotationProcessorOption) map[string]any {
	var m map[string]any
	return m
}
