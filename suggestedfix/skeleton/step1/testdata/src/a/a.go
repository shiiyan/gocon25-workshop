package a

func example1(x interface{}) interface{} { // want "interface found"
	return x
}

func example2(x any) any {
	return x
}

type MyStruct struct {
	Field1 interface{} // want "interface{} found"
	Field2 any
}

var globalVar interface{} = "hello" // want "interface{} found"

func example3() {
	var local interface{} = 42 // want "interface{} found"
	_ = local
}

type GenericType[T any] struct {
	Value T
}

type OldGenericType[T interface{}] struct { // want "interface{} found"
	Value T
}

func example4(items []interface{}) { // want "interface{} found"
	for _, item := range items {
		_ = item
	}
}

func example5() map[string]interface{} { // want "interface{} found"
	return map[string]interface{}{ // want "interface{} found"
		"key": "value",
	}
}

func example6(ch chan interface{}) { // want "interface{} found"
	ch <- "test"
}
