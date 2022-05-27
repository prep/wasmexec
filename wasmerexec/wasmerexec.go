package wasmerexec

import (
	"github.com/prep/wasmexec"

	"github.com/wasmerio/wasmer-go/wasmer"
)

// Import the Go JavaScript functions.
func Import(store *wasmer.Store, instance wasmexec.Instance) *wasmer.ImportObject {
	mod := wasmexec.NewModuleGo(instance)

	wasmExit := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.WasmExit(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	wasmWrite := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.WasmWrite(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	resetMemoryDataView := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ResetMemoryDataView(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	nanotime1 := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.Nanotime1(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	walltime := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.Walltime(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	scheduleTimeoutEvent := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ScheduleTimeoutEvent(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	clearTimeoutEvent := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ClearTimeoutEvent(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	getRandomData := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.GetRandomData(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	finalizeRef := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.FinalizeRef(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	stringVal := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.StringVal(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueGet := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueGet(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueSet := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueSet(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueDelete := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueDelete(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueIndex := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueIndex(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueSetIndex := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueSetIndex(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueCall := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueCall(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueInvoke := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueInvoke(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueNew := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueNew(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueLength := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueLength(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valuePrepareString := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValuePrepareString(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueLoadString := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueLoadString(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	valueInstanceOf := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.ValueInstanceOf(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	copyBytesToGo := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.CopyBytesToGo(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	copyBytesToJS := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.CopyBytesToJS(args[0].I32())
			return []wasmer.Value{}, nil
		},
	)

	debug := wasmer.NewFunction(
		store,
		wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			mod.Debug(args[0].I32())
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

	return imports
}
