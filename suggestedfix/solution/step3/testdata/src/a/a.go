package a

// Empty interface declarations should trigger diagnostics.
type Empty interface{} // want "interface{} can be replaced with any"

type EmptyAlias = interface{} // want "interface{} can be replaced with any"

// Non-empty interfaces (with methods or embedded interfaces) should be ignored.
type NonEmpty interface {
	Method()
}

type Embedded interface {
	error
}

func AcceptEmpty(x interface{}) { // want "interface{} can be replaced with any"
	_ = x
}

func AcceptNonEmpty(x interface{ Close() error }) {
	_ = x
}

var Value interface{} = struct{}{} // want "interface{} can be replaced with any"

var NonEmptyValue interface{ Read([]byte) (int, error) } = nil

// Type block with both empty and non-empty interfaces.
type (
	AnotherEmpty    interface{} // want "interface{} can be replaced with any"
	AnotherNonEmpty interface {
		Write([]byte) (int, error)
	}
)

func AcceptMap(m map[interface{}]string) { // want "interface{} can be replaced with any"
	_ = m
}

type Wrapper[T interface{}] struct { // want "interface{} can be replaced with any"
	value T
}

func AcceptGeneric[T interface{}](v T) { // want "interface{} can be replaced with any"
	_ = v
}
