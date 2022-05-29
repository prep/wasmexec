package wasmexec

// Instance describes an instance of a Wasm module.
type Instance interface {
	Memory

	// Debug and Error are used for development purposes.
	Debug(format string, params ...interface{})
	Error(format string, params ...interface{})

	// Get the SP value.
	GetSP() (int32, error)

	// Resume the execution of Go code until it needs to wait for an event.
	Resume() error

	// Write is called whenever the program wants to write to a file descriptor.
	Write(fd int, b []byte) (int, error)

	// Exit is called whenever this app closes.
	Exit(code int)
}
