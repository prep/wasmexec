module github.com/prep/wasmexec/wazeroexec

go 1.18

replace github.com/prep/wasmexec => ../

require (
	github.com/prep/wasmexec main
	github.com/tetratelabs/wazero v0.0.0-20220523092326-5ed31d3c495d
)
