# wapc
This is the implementation of the [waPC](https://github.com/wapc) interface for WebAssembly guest modules written in Go. It allows a `wasmexec` host to invoke procedures inside a Go compiled guest and similarly for the guest to invoke procedures exposed by the `wasmexec` host.

This implementation is built on top of the [js.FuncOf()](https://pkg.go.dev/syscall/js#FuncOf) functionality and thus does not leverage the imports and exports dictated by the waPC standard. This means this is not a real waPC implementation and can not be used by a waPC-compliant host that has not implemented this `wasmexec` package.

## Host
On the host instance you can call `Invoke()` on `*wasmexec.Module` to send something to the guest:

```go
result, err := mod.Invoke(ctx, "hello", []byte(`Hello World`))
```

The host can also receive events by implementing `HostCall()` on the instance:

```go
func (instance *Instance) HostCall(binding, namespace, operation string, payload []byte) ([]byte, error) {
    // ...
}
```

For examples of host implementations, check the `example` in each runtime-specific directory.

## Guest
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

The [waPC example](../examples/wapc) shows a simple guest implementation that gets called by the runtime examples.