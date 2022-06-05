package wasmexec

// jsFunction describes the constructor of an jsObject.
type jsFunction struct {
	name string
	fn   func(args []any) any
}

// newjsFunction returns a new function.
func newjsFunction(fn func(args []any) any) *jsFunction {
	return &jsFunction{fn: fn}
}

// Name returns the name of the constructor type.
func (fn jsFunction) Name() string {
	return fn.name
}

// jsProperties describe the properties on an object. This can either be a
// function or a value.
type jsProperties map[string]any

// jsObject describes a JSON object.
type jsObject struct {
	properties jsProperties
}

// jsArray describes an array of elements.
type jsArray struct {
	elements []any
}

// jsUint8Array describes a byte slice.
type jsUint8Array struct {
	data []byte
}

// jsString represents a stored string.
type jsString struct {
	data string
}
