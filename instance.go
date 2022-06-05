package wasmexec

// Instance describes an instance of a Wasm module.
type Instance interface {
	Memory

	// Get the SP value.
	GetSP() (uint32, error)

	// Resume the execution of Go code until it needs to wait for an event.
	Resume() error

	// Write is called whenever the program wants to write to a file descriptor.
	Write(fd int, b []byte) (int, error)
}
