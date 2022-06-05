package wazeroexec

import (
	"context"

	"github.com/prep/wasmexec"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Import the Go JavaScript functions.
func Import(ctx context.Context, runtime wazero.Runtime, ns wazero.Namespace, instance wasmexec.Instance) (*wasmexec.Module, error) {
	mod := wasmexec.New(instance)

	funcs := map[string]any{
		"runtime.wasmExit":              mod.WasmExit,
		"runtime.wasmWrite":             mod.WasmWrite,
		"runtime.resetMemoryDataView":   mod.ResetMemoryDataView,
		"runtime.nanotime1":             mod.Nanotime1,
		"runtime.walltime":              mod.Walltime,
		"runtime.scheduleTimeoutEvent":  mod.ScheduleTimeoutEvent,
		"runtime.clearTimeoutEvent":     mod.ClearTimeoutEvent,
		"runtime.getRandomData":         mod.GetRandomData,
		"syscall/js.finalizeRef":        mod.FinalizeRef,
		"syscall/js.stringVal":          mod.StringVal,
		"syscall/js.valueGet":           mod.ValueGet,
		"syscall/js.valueSet":           mod.ValueSet,
		"syscall/js.valueDelete":        mod.ValueDelete,
		"syscall/js.valueIndex":         mod.ValueIndex,
		"syscall/js.valueSetIndex":      mod.ValueSetIndex,
		"syscall/js.valueCall":          mod.ValueCall,
		"syscall/js.valueInvoke":        mod.ValueInvoke,
		"syscall/js.valueNew":           mod.ValueNew,
		"syscall/js.valueLength":        mod.ValueLength,
		"syscall/js.valuePrepareString": mod.ValuePrepareString,
		"syscall/js.valueLoadString":    mod.ValueLoadString,
		"syscall/js.valueInstanceOf":    mod.ValueInstanceOf,
		"syscall/js.copyBytesToGo":      mod.CopyBytesToGo,
		"syscall/js.copyBytesToJS":      mod.CopyBytesToJS,
		"debug":                         mod.Debug,
	}

	_, err := runtime.NewModuleBuilder("go").ExportFunctions(funcs).Instantiate(ctx, ns)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

// Memory wraps a wazero Memory module and slaps on top of it the getter and
// setter functions that the wasmexec runtime needs.
type Memory struct {
	api.Memory
}

// NewMemory returns a new Memory.
func NewMemory(mem api.Memory) *Memory {
	return &Memory{Memory: mem}
}

func (mem *Memory) GetUInt32(offset uint32) (uint32, error) {
	val, ok := mem.ReadUint32Le(context.Background(), offset)
	if !ok {
		return 0, wasmexec.ErrFault
	}

	return val, nil
}

func (mem *Memory) GetInt64(offset uint32) (int64, error) {
	val, ok := mem.ReadUint64Le(context.Background(), offset)
	if !ok {
		return 0, wasmexec.ErrFault
	}

	return int64(val), nil
}

func (mem *Memory) GetFloat64(offset uint32) (float64, error) {
	val, ok := mem.ReadFloat64Le(context.Background(), offset)
	if !ok {
		return 0, wasmexec.ErrFault
	}

	return val, nil
}

func (mem *Memory) Mem(offset, length uint32) ([]byte, error) {
	data, ok := mem.Read(context.Background(), offset, length)
	if !ok {
		return nil, wasmexec.ErrFault
	}

	return data, nil
}

func (mem *Memory) SetUInt8(offset uint32, val uint8) error {
	ok := mem.WriteByte(context.Background(), offset, val)
	if !ok {
		return wasmexec.ErrFault
	}

	return nil
}

func (mem *Memory) SetUInt32(offset, val uint32) error {
	ok := mem.WriteUint32Le(context.Background(), offset, val)
	if !ok {
		return wasmexec.ErrFault
	}

	return nil
}

func (mem *Memory) SetInt64(offset uint32, val int64) error {
	ok := mem.WriteUint64Le(context.Background(), offset, uint64(val))
	if !ok {
		return wasmexec.ErrFault
	}

	return nil
}

func (mem *Memory) SetFloat64(offset uint32, val float64) error {
	ok := mem.WriteFloat64Le(context.Background(), offset, val)
	if !ok {
		return wasmexec.ErrFault
	}

	return nil
}
