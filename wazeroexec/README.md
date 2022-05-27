# wazeroexec
This package provides an import hook for [wazero](https://github.com/tetratelabs/wazero). See the [example](example/) directory for a working implementation.

Usage:

```go
import (
  // ...
  "github.com/prep/wasmexec/wazeroexec"
  // ...
)

func main() {
  // ...
  if err = wazeroexec.Import(ctx, runtime, instance); err != nil {
    // handle error
  }
  // ...
}
```