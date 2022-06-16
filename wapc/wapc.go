//go:build js && wasm

package wapc

import (
	"errors"
	"syscall/js"
)

var (
	jswaPC     = js.Global().Get("wapc")
	array      = js.Global().Get("Array")
	uint8Array = js.Global().Get("Uint8Array")
)

func init() {
	js.Global().Set("__guest_call", js.FuncOf(guestCall))
}

type Function func(payload []byte) ([]byte, error)

type Functions map[string]Function

var allFunctions = Functions{}

func RegisterFunction(name string, fn Function) {
	allFunctions[name] = fn
}

func RegisterFunctions(functions Functions) {
	for name, fn := range functions {
		RegisterFunction(name, fn)
	}
}

// HostCall invokes an operation on the host. The host uses namespace and
// operation to route the payload to the appropriate operation. The host will
// return a response payload if successful.
func HostCall(binding, namespace, operation string, payload []byte) ([]byte, error) {
	result := jswaPC.Call("__host_call", stringToJS(binding), stringToJS(namespace), stringToJS(operation), bytesToJS(payload))

	// No validation of the result here as that is kind of expensive.

	var response []byte
	if respValue := result.Index(0); !respValue.IsNull() {
		response = bytesFromJS(respValue)
	}

	var err error
	if errValue := result.Index(1); !errValue.IsNull() {
		err = errors.New(errValue.String())
	}

	return response, err
}

func guestCall(_ js.Value, args []js.Value) any {
	switch {
	// Make sure there are 2 arguments.
	case len(args) != 2:
		return false
	// Make sure the 1st one is a string.
	case args[0].Type() != js.TypeString:
		return false
	// Make sure the 2nd one is an object.
	case args[1].Type() != js.TypeObject:
		return false
	// Make sure the 2nd one is derived from Uint8Array, meaning a []byte.
	case !args[1].InstanceOf(uint8Array):
		return false
	}

	// Get the operation.
	operation := args[0].String()

	// Copy the payload over from the host to this guest.
	payload := bytesFromJS(args[1])

	// Find the function that matches the operation name.
	fn, ok := allFunctions[operation]
	if !ok {
		guestError(`Could not find function "` + operation + `"`)
		return false
	}

	// Call the operation function.
	response, err := fn(payload)
	if err != nil {
		guestError(err.Error())
		return false
	}

	guestResponse(response)
	return true
}

// guestResponse sets the guest response.
func guestResponse(payload []byte) {
	jswaPC.Call("__guest_response", bytesToJS(payload))
}

// guestError sets the guest error.
func guestError(message string) {
	jswaPC.Call("__guest_error", stringToJS(message))
}

// bytesFromJS converts a js.Value to a []byte.
func bytesFromJS(v js.Value) []byte {
	s := make([]byte, v.Length())
	js.CopyBytesToGo(s, v)
	return s
}

// bytesToJS converts a []byte to a js.Value.
func bytesToJS(d []byte) js.Value {
	a := uint8Array.New(len(d))
	js.CopyBytesToJS(a, d)
	return a
}

// stringToJS converts a string to a js.Value.
func stringToJS(s string) js.Value {
	return bytesToJS([]byte(s))
}
