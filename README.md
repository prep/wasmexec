# wasmexec
wasmexec is runtime-agnostic implementation of Go's [wasm_exec.js](https://github.com/golang/go/blob/master/misc/wasm/wasm_exec.js) in Go. It currently has import hooks for [wasmer](wasmerexec/), [wasmtime](wasmtimexec/) and [wazero](wazeroexec/). Each runtime-dedicated package has its own example of an implementation that can run any of the [examples](examples/).

This implementation was made possible by allowing me to peek at mattn's [implementation](https://github.com/mattn/gowasmer/) as well as Vedhavyas Singareddi's [go-wasm-adapter](https://github.com/go-wasm-adapter/go-wasm/).

**NOTE:** This implementation is still highly experimental. Please let me know what breaks in an issue with an example piece of code.