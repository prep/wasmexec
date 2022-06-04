# wasmexec
wasmexec is runtime-agnostic implementation of Go's [wasm_exec.js](https://github.com/golang/go/blob/master/misc/wasm/wasm_exec.js) in Go. It currently has import hooks for [wasmer](wasmerexec/), [wasmtime](wasmtimexec/) and [wazero](wazeroexec/). Each runtime-dedicated package has its own example of an implementation that can run any of the [examples](examples/).

## js.FuncOf()
The guest can use [js.FuncOf()](https://pkg.go.dev/syscall/js#FuncOf) to create functions that can be called on the host via `Call()` on `*wasmexec.ModuleGo`.

```go
var uint8Array = js.Global().Get("Uint8Array")

func main() {
    myEvent := js.FuncOf(func(this js.Value, args []js.Value) any {
        arg := args[0]

        if arg.InstanceOf(uint8Array) {
            dst := make([]byte, arg.Length())
            js.CopyBytesToGo(dst, arg)
						
            fmt.Printf("Received: %v\n", string(dst))
        }
				
        return nil
    }
		
    js.Global().Set("myEvent", myEvent)
}
```

On the host these functions can be called using `Call()` on `*wasmexec.ModuleGo`:

```go
mod.Call("myEvent", []byte("Hello World!"))

```

## waPC
wasmexec supports a fake waPC implementation built on top of the above-mentioned `js.FuncOf()` functionality. Due to the fact that the current Go Wasm compiler cannot export functions nor require import functions, a workaround was created to get the same effect.

On the host instance you can call `Invoke()` on the `*wasmexec.ModuleGo` type to send something to the guest:

```go
result, err := mod.Invoke(context.Background(), "hello", []byte(`Hello World`))
```

The host can also receive events by implementing `HostCall()` on the instance:

```go
func (instance *Instance) HostCall(binding, namespace, operation string, payload []byte) ([]byte, error) {
    // ...
}
```

On the guest, events from the host's `Invoke()` calls can be received as follows:

```go
import (
    "fmt"

    "github.com/prep/wasmexec/wapc"
)

func hello(payload []byte) ([]byte, error) {
    return []byte("Hello back!"), nil
}

func main() {
    wapc.RegisterFunctions(wapc.Functions{
        "hello": hello,
    })
}
```

The guest can also send events to the host, which will be received by the host's `HostCall()` function:

```go
resp, err := wapc.HostCall("myBinding", "sample", "hello", []byte("Guest"))

```

## Acknowledgements
This implementation was made possible by allowing me to peek at mattn's [implementation](https://github.com/mattn/gowasmer/) as well as Vedhavyas Singareddi's [go-wasm-adapter](https://github.com/go-wasm-adapter/go-wasm/).
