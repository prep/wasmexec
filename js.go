package wasmexec

// errno describes an error "number".
type errno string

// This errno list is a subset of the errors in syscall/tables_js.go.
const (
	eNOSYS errno = "ENOSYS"
)

// errorResponse returns a errno callback response.
func errorResponse(code errno) []any {
	return []any{jsProperties{"code": string(code)}}
}

// errCallback calls the callback with the specified errno code.
func errorCallback(code errno) *jsFunction {
	return &jsFunction{
		fn: func(args []any) any {
			if len(args) == 0 {
				return nil
			}

			// The last item in the list should be a callback function, according to
			// the fsCall() function in syscall/fs.js.go.
			callback, ok := args[len(args)-1].(*jsFunction)
			if !ok {
				return nil
			}

			callback.fn(errorResponse(code))
			return nil
		},
	}
}

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
