package wasmtimexec

import (
	"github.com/prep/wasmexec"

	"github.com/bytecodealliance/wasmtime-go"
)

var i32 = wasmtime.NewValType(wasmtime.KindI32)

// Import the Go JavaScript functions.
func Import(store *wasmtime.Store, linker *wasmtime.Linker, instance wasmexec.Instance) (*wasmexec.ModuleGo, error) {
	var err error
	define := func(module, name string, item wasmtime.AsExtern) {
		if err != nil {
			return
		}

		err = linker.Define(module, name, item)
	}

	mod := wasmexec.NewModuleGo(instance)

	define("go", "runtime.wasmExit", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.WasmExit(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.wasmWrite", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.WasmWrite(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.resetMemoryDataView", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ResetMemoryDataView(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.nanotime1", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.Nanotime1(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.walltime", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.Walltime(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.scheduleTimeoutEvent", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ScheduleTimeoutEvent(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.clearTimeoutEvent", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ClearTimeoutEvent(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.getRandomData", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.GetRandomData(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.finalizeRef", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.FinalizeRef(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.stringVal", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.StringVal(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueGet", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueGet(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueSet", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueSet(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueDelete", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueDelete(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueIndex", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueIndex(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueSetIndex", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueSetIndex(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueCall", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueCall(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueInvoke", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueCall(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueNew", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueNew(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueLength", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueLength(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valuePrepareString", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValuePrepareString(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueLoadString", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueLoadString(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueInstanceOf", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueInstanceOf(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.copyBytesToGo", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.CopyBytesToGo(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.copyBytesToJS", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.CopyBytesToJS(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "debug", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.Debug(uint32(args[0].I32()))
			return []wasmtime.Val{}, nil
		}),
	)

	return mod, err
}
