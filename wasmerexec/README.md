# wasmerexec
This package provides an import hook for [wasmer-go](https://github.com/wasmerio/wasmer-go). See the [example](example/) directory for a working implementation.

Usage:

```go
import "github.com/prep/wasmexec/wasmerexec"
```

```go
imports := wasmerexec.Import(store, instance)
```