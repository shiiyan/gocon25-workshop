package a

func example1(x interface{}) interface{} { // want "interface{} can be replaced with any" "interface{} can be replaced with any"
	return x
}

func example2(x any) any {
	return x
}

type MyStruct struct {
	Field1 interface{} // want "interface{} can be replaced with any"
	Field2 any
}

var globalVar interface{} = "hello" // want "interface{} can be replaced with any"

func example3() {
	var local interface{} = 42 // want "interface{} can be replaced with any"
	_ = local
}

type GenericType[T any] struct {
	Value T
}

type OldGenericType[T interface{}] struct { // want "interface{} can be replaced with any"
	Value T
}

func example4(items []interface{}) { // want "interface{} can be replaced with any"
	for _, item := range items {
		_ = item
	}
}

func example5() map[string]interface{} { // want "interface{} can be replaced with any"
	return map[string]interface{}{ // want "interface{} can be replaced with any"
		"key": "value",
	}
}

func example6(ch chan interface{}) { // want "interface{} can be replaced with any"
	ch <- "test"
}