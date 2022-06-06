# wasmexec
wasmexec is runtime-agnostic implementation of Go's [wasm_exec.js](https://github.com/golang/go/blob/master/misc/wasm/wasm_exec.js) in Go. It currently has import hooks for [wasmer](wasmerexec/), [wasmtime](wasmtimexec/) and [wazero](wazeroexec/). Each runtime-dedicated package has its own example of an implementation that can run any of the [examples](examples/).

## 1. Minimum implementation
When a Go Wasm module is instantiated, it needs to be wrapped up in a structure that implements the minimum set of methods in order to run. `wasmexec.New()` accepts the [Instance](instance.go) interface.


```go
type Instance interface {
    Memory

    GetSP() (uint32, error)
    Resume() error
}
```

The `GetSP()` and `Resume()` methods are calls directly to the Go Wasm exports. The [Memory](memory.go) interface wraps the Go Wasm module's memory.

```go
type Memory interface {
    Range(offset, length uint32) ([]byte, error)
		
    GetUInt32(offset uint32) (uint32, error)
    GetInt64(offset uint32) (int64, error)
    GetFloat64(offset uint32) (float64, error)
    SetUInt8(offset uint32, val uint8) error
    SetUInt32(offset, val uint32) error
    SetInt64(offset uint32, val int64) error
    SetFloat64(offset uint32, val float64) error
}
```

If your runtime exposes the memory as a `[]byte` (as wasmer and wasmtime do) then you can easily use the `NewMemory()` function to satisfy this interface. If not, a custom implementation needs to be written (like wazero).

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
If the `errorLogger` interface is implemented, `Error()` is called for any error that might pop up during execution. At this stage, it is probably useful to implement this as this package isn't battle tested yet.

```go
type errorLogger interface {
    Error(format string, params ...any)
}
```

### 2.3. Writer
If the `fdWriter` interface is implemented, `Write()` is called for any data being sent to `stdout` or `stderr`. It is highly recommended that this is implemented.

```go
type fdWriter interface {
    Write(fd int, data []byte) (n int, err error)
}

```

### 2.4. Exiting
If the `exiter` interface is implemented, `Exit()` is called whenever the call to the `run()` Wasm function is done.

```go
type exiter interface {
    Exit(code int)
}
```

### 2.5. waPC
If the `hostCaller` interface is implemented, `HostCall()` is called whenever the waPC guest sends information to the host. More information on this can be found in the [wapc](wapc/) package.

```go
type hostCaller interface {
    HostCall(string, string, string, []byte) ([]byte, error)
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

## 4. Acknowledgements
This implementation was made possible by allowing me to peek at mattn's [implementation](https://github.com/mattn/gowasmer/) as well as Vedhavyas Singareddi's [go-wasm-adapter](https://github.com/go-wasm-adapter/go-wasm/).
