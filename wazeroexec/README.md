# wazeroexec
This package provides an import hook for [wazero](https://github.com/tetratelabs/wazero). See the [example](example/) directory for a working implementation.

Usage:

```go
import "github.com/prep/wasmexec/wazeroexec"
```

```go
err := wazeroexec.Import(ctx, runtime, instance)
```

Or if you want to import these functions on your namespace:

```go
err := wazeroexec.ImportWithNamespace(ctx, runtime, ns, instance)
```