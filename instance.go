package wasmexec

// Instance describes an instance of a Wasm module.
type Instance interface {
	Memory

	// Get the SP value.
	GetSP() (uint32, error)

	// Resume the execution of Go code until it needs to wait for an event.
	Resume() error
}
