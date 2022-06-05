package wasmexec

// Function describes the constructor of an Object.
type Function struct {
	name string
	fn   func(args []any) any
}

// Name returns the name of the constructor type.
func (fn Function) Name() string {
	return fn.name
}

// newFuncObject returns a new function.
func newFuncObject(fn func(args []any) any) *Function {
	return &Function{fn: fn}
}

// Properties describe the properties on an object. This can either be a
// function or a value.
type Properties map[string]any

// Object describes a JSON object.
type Object struct {
	properties Properties
}

// Array describes an array of elements.
type Array struct {
	elements []any
}

// Uint8Array describes a byte slice.
type Uint8Array struct {
	data []byte
}

// String represents a stored string.
type String struct {
	data string
}
