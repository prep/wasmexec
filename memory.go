package wasmexec

import (
	"encoding/binary"
	"errors"
	"math"
)

// ErrFault is returned whenever memory was accessed that was not addressable.
var ErrFault = errors.New("bad address")

type Memory interface {
	Mem(offset, length uint32) ([]byte, error)
	GetUInt32(offset uint32) (uint32, error)
	GetInt64(offset uint32) (int64, error)
	GetFloat64(offset uint32) (float64, error)
	SetUInt8(offset uint32, val uint8) error
	SetUInt32(offset, val uint32) error
	SetInt64(offset uint32, val int64) error
	SetFloat64(offset uint32, val float64) error
}

type memory []byte

// NewMemory returns a new Memory.
func NewMemory(mem []byte) Memory {
	return memory(mem)
}

// Mem returns part of memory block.
func (mem memory) Mem(offset, length uint32) ([]byte, error) {
	if int(offset+length) >= len(mem) {
		return nil, ErrFault
	}

	return mem[offset : offset+length], nil
}

// GetUInt32 returns an uint32 value.
func (mem memory) GetUInt32(offset uint32) (uint32, error) {
	if int(offset+4) >= len(mem) {
		return 0, ErrFault
	}

	return binary.LittleEndian.Uint32(mem[offset:]), nil
}

// GetInt64 returns an int64 value.
func (mem memory) GetInt64(offset uint32) (int64, error) {
	if int(offset+8) >= len(mem) {
		return 0, ErrFault
	}

	return int64(binary.LittleEndian.Uint64(mem[offset:])), nil
}

// GetFloat64 returns a float64 value.
func (mem memory) GetFloat64(offset uint32) (float64, error) {
	if int(offset+8) >= len(mem) {
		return 0, ErrFault
	}

	return math.Float64frombits(binary.LittleEndian.Uint64(mem[offset:])), nil
}

// SetUInt8 sets an uint8 value.
func (mem memory) SetUInt8(offset uint32, val uint8) error {
	if int(offset+1) >= len(mem) {
		return ErrFault
	}

	mem[offset] = val
	return nil
}

// SetUInt32 sets an uint32 value.
func (mem memory) SetUInt32(offset, val uint32) error {
	if int(offset+4) >= len(mem) {
		return ErrFault
	}

	binary.LittleEndian.PutUint32(mem[offset:], val)
	return nil
}

// SetInt64 sets an int64 value.
func (mem memory) SetInt64(offset uint32, val int64) error {
	if int(offset+8) >= len(mem) {
		return ErrFault
	}

	binary.LittleEndian.PutUint64(mem[offset:], uint64(val))
	return nil
}

// SetFloat64 sets a float64 value.
func (mem memory) SetFloat64(offset uint32, val float64) error {
	if int(offset+8) >= len(mem) {
		return ErrFault
	}

	binary.LittleEndian.PutUint64(mem[offset:], math.Float64bits(val))
	return nil
}
