# wasmtimexec
This package provides an import hook for [wasmtime-go](https://github.com/bytecodealliance/wasmtime-go). See the [example](example/) directory for a working implementation.

Usage:

```go
import "github.com/prep/wasmexec/wasmtimexec"
```

```go
err := wasmtimexec.Import(store, linker, instance)
```