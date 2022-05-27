package wasmexec

// Instance describes an instance of a Wasm module.
type Instance interface {
	// Debug and Error are used for development purposes.
	Debug(format string, params ...interface{})
	Error(format string, params ...interface{})

	// These methods are implemented by the Memory type which should be composed
	// into an Instance type.
	Mem(offset, length int32) ([]byte, error)
	GetUInt32(offset int32) (uint32, error)
	GetInt64(offset int32) (int64, error)
	GetFloat64(offset int32) (float64, error)
	SetUInt8(offset int32, val uint8) error
	SetUInt32(offset int32, val uint32) error
	SetInt64(offset int32, val int64) error
	SetFloat64(offset int32, val float64) error

	// Get the SP value.
	GetSP() (int32, error)

	// Resume the execution of Go code until it needs to wait for an event.
	Resume() error

	// Write is called whenever the program wants to write to a file descriptor.
	Write(fd int, b []byte) (int, error)

	// Exit is called whenever this app closes.
	Exit(code int)
}
