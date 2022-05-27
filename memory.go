package wasmexec

import (
	"encoding/binary"
	"errors"
	"math"
)

// ErrFault is returned whenever memory was accessed that was not addressable.
var ErrFault = errors.New("bad address")

// Memory describes a block of memory used by a Wasm instance.
type Memory []byte

// Mem returns part of memory block.
func (mem Memory) Mem(offset, length int32) ([]byte, error) {
	if int(offset+length) >= len(mem) {
		return nil, ErrFault
	}

	return mem[offset : offset+length], nil
}

// GetUInt32 returns an uint32 value.
func (mem Memory) GetUInt32(offset int32) (uint32, error) {
	if int(offset+4) >= len(mem) {
		return 0, ErrFault
	}

	return binary.LittleEndian.Uint32(mem[offset:]), nil
}

// GetInt64 returns an int64 value.
func (mem Memory) GetInt64(offset int32) (int64, error) {
	if int(offset+8) >= len(mem) {
		return 0, ErrFault
	}

	return int64(binary.LittleEndian.Uint64(mem[offset:])), nil
}

// GetFloat64 returns a float64 value.
func (mem Memory) GetFloat64(offset int32) (float64, error) {
	if int(offset+8) >= len(mem) {
		return 0, ErrFault
	}

	return math.Float64frombits(binary.LittleEndian.Uint64(mem[offset:])), nil
}

// SetUInt8 sets an uint8 value.
func (mem Memory) SetUInt8(offset int32, val uint8) error {
	if int(offset+1) >= len(mem) {
		return ErrFault
	}

	mem[offset] = val
	return nil
}

// SetInt32 sets an int32 value.
func (mem Memory) SetInt32(offset, val int32) error {
	return mem.SetUInt32(offset, uint32(val))
}

// SetUInt32 sets an uint32 value.
func (mem Memory) SetUInt32(offset int32, val uint32) error {
	if int(offset+4) >= len(mem) {
		return ErrFault
	}

	binary.LittleEndian.PutUint32(mem[offset:], val)
	return nil
}

// SetInt64 sets an int64 value.
func (mem Memory) SetInt64(offset int32, val int64) error {
	if int(offset+8) >= len(mem) {
		return ErrFault
	}

	binary.LittleEndian.PutUint64(mem[offset:], uint64(val))
	return nil
}

// SetFloat64 sets a float64 value.
func (mem Memory) SetFloat64(offset int32, val float64) error {
	if int(offset+8) >= len(mem) {
		return ErrFault
	}

	binary.LittleEndian.PutUint64(mem[offset:], math.Float64bits(val))
	return nil
}
