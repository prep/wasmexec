package wasmerexec

import (
	"github.com/prep/wasmexec"

	"github.com/wasmerio/wasmer-go/wasmer"
)

// Import the Go JavaScript functions.
func Import(store *wasmer.Store, instance wasmexec.Instance) (*wasmer.ImportObject, *wasmexec.Module) {
	mod := wasmexec.New(instance)

	wasmExit := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.WasmExit(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	wasmWrite := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.WasmWrite(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	resetMemoryDataView := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ResetMemoryDataView(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	nanotime1 := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.Nanotime1(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	walltime := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.Walltime(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	scheduleTimeoutEvent := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ScheduleTimeoutEvent(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	clearTimeoutEvent := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ClearTimeoutEvent(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	getRandomData := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.GetRandomData(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	finalizeRef := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.FinalizeRef(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	stringVal := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.StringVal(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueGet := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueGet(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueSet := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueSet(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueDelete := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueDelete(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueIndex := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueIndex(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueSetIndex := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueSetIndex(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueCall := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueCall(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueInvoke := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueInvoke(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueNew := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueNew(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueLength := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueLength(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valuePrepareString := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValuePrepareString(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueLoadString := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueLoadString(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	valueInstanceOf := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueInstanceOf(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	copyBytesToGo := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.CopyBytesToGo(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	copyBytesToJS := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.CopyBytesToJS(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	debug := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.Debug(uint32(args[0].I32()))
			return []wasmer.Value{}, nil
		},
	)

	imports := wasmer.NewImportObject()
	imports.Register("go", map[string]wasmer.IntoExtern{
		"runtime.wasmExit":              wasmExit,
		"runtime.wasmWrite":             wasmWrite,
		"runtime.resetMemoryDataView":   resetMemoryDataView,
		"runtime.nanotime1":             nanotime1,
		"runtime.walltime":              walltime,
		"runtime.scheduleTimeoutEvent":  scheduleTimeoutEvent,
		"runtime.clearTimeoutEvent":     clearTimeoutEvent,
		"runtime.getRandomData":         getRandomData,
		"syscall/js.finalizeRef":        finalizeRef,
		"syscall/js.stringVal":          stringVal,
		"syscall/js.valueGet":           valueGet,
		"syscall/js.valueSet":           valueSet,
		"syscall/js.valueDelete":        valueDelete,
		"syscall/js.valueIndex":         valueIndex,
		"syscall/js.valueSetIndex":      valueSetIndex,
		"syscall/js.valueCall":          valueCall,
		"syscall/js.valueInvoke":        valueInvoke,
		"syscall/js.valueNew":           valueNew,
		"syscall/js.valueLength":        valueLength,
		"syscall/js.valuePrepareString": valuePrepareString,
		"syscall/js.valueLoadString":    valueLoadString,
		"syscall/js.valueInstanceOf":    valueInstanceOf,
		"syscall/js.copyBytesToGo":      copyBytesToGo,
		"syscall/js.copyBytesToJS":      copyBytesToJS,
		"debug":                         debug,
	})

	return imports, mod
}
