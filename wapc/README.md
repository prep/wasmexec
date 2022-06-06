# wapc
This is the implementation of the **waPC** API (not standard) for WebAssembly guest modules written in Go. It allows a `wasmexec` host to invoke procedures inside a Go compiled guest and similarly for the guest to invoke procedures exposed by the `wsamexec` host.

This implementation is built on top of the [js.FuncOf()](https://pkg.go.dev/syscall/js#FuncOf) functionality and thus does not leverage the standard imports and exports dictated by the waPC standard. This means this is not a real waPC implementation and can not be used by a waPC-compliant host that has not implemented this `wasmexec` package.

## Usage
The [waPC example](../examples/wapc) shows a simple implementation that gets called by the runtime examples.