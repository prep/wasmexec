# examples
This directory contains some Wasm examples. Compile them as follows:

```
env GOOS=js GOARCH=wasm go build -o example1.wasm example1.go
env GOOS=js GOARCH=wasm go build -o example2.wasm example2.go
```