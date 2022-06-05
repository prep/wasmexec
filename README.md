# wasmexec
wasmexec is runtime-agnostic implementation of Go's [wasm_exec.js](https://github.com/golang/go/blob/master/misc/wasm/wasm_exec.js) in Go. It currently has import hooks for [wasmer](wasmerexec/), [wasmtime](wasmtimexec/) and [wazero](wazeroexec/). Each runtime-dedicated package has its own example of an implementation that can run any of the [examples](examples/).

## 1. Basic implementation
The per-runtime examples are a good starter, but you basically instantiate a Go Wasm module and wrap that instance up in a custom struct that implements several methods. At a minimum, your wrapper struct needs to satisfiy the [Instance](instance.go) interface.

```go
type Instance interface {
    Memory

    GetSP() (uint32, error)
    Resume() error
    Write(fd int, b []byte) (int, error)
}
```

The `GetSP()` and `Resume()` methods are calls directly to the Go Wasm exports. The `Write()` method is called for writes to `stdout` or `stderr`.

The [Memory](memory.go) interface describes what is needed to read and write to the Go Wasm memory. If your runtime exposes the memory as a `[]byte` (as wasmer and wasmtime do) then you can easily use the `NewMemory()` function to satisfy this interface.

```go
type Memory interface {
    Mem(offset, length uint32) ([]byte, error)
		
    GetUInt32(offset uint32) (uint32, error)
    GetInt64(offset uint32) (int64, error)
    GetFloat64(offset uint32) (float64, error)
    SetUInt8(offset uint32, val uint8) error
    SetUInt32(offset, val uint32) error
    SetInt64(offset uint32, val int64) error
    SetFloat64(offset uint32, val float64) error
}
```

## 2. Optional implementation
Your wrapper struct can also implement additional methods that are called when applicable.

### 2.1. Debug logging
If the `debugLogger` interface is implemented, `Debug()` is called with debug messages. This is only useful when you're debugging issues with this package.

```go
type debugLogger interface {
  Debug(format string, params ...any)
}
```

### 2.2. Error logging
If the `errorLogger` interface is implement, `Error()` is called for any error that might pop up during execution. At this stage, it is probably useful to implement this as this package isn't battle tested yet.

```go
type errorLogger interface {
  Error(format string, params ...any)
}
```

### 2.3. Exiting
If the `exiter` interface is implemented, `Exit()` is called whenver the call to the `run()` Wasm function is done.

```go
type exiter interface {
  Exit(code int)
}
```

## 3. js.FuncOf()
The guest can use [js.FuncOf()](https://pkg.go.dev/syscall/js#FuncOf) to create functions that can be called from the host.

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

On the host these functions can be called using `Call()` on `*wasmexec.Module`:

```go
mod.Call("myEvent", []byte("Hello World!"))
```

## 4. waPC
wasmexec supports a fake waPC implementation built on top of the above-mentioned `js.FuncOf()` functionality. Due to the fact that the current Go Wasm compiler cannot export functions nor require import functions, a workaround was created to get the same effect.

On the host instance you can call `Invoke()` on the `*wasmexec.Module` type to send something to the guest:

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

## 5. Acknowledgements
This implementation was made possible by allowing me to peek at mattn's [implementation](https://github.com/mattn/gowasmer/) as well as Vedhavyas Singareddi's [go-wasm-adapter](https://github.com/go-wasm-adapter/go-wasm/).
