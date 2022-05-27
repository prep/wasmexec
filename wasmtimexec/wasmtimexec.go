package wasmtimexec

import (
	"github.com/prep/wasmexec"

	"github.com/bytecodealliance/wasmtime-go"
)

var (
	i32 = wasmtime.NewValType(wasmtime.KindI32)
	i64 = wasmtime.NewValType(wasmtime.KindI64)
)

// Import the Go JavaScript functions.
func Import(store *wasmtime.Store, linker *wasmtime.Linker, instance wasmexec.Instance) error {
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
			mod.WasmExit(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.wasmWrite", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.WasmWrite(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.resetMemoryDataView", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ResetMemoryDataView(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.nanotime1", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.Nanotime1(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.walltime", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.Walltime(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.scheduleTimeoutEvent", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ScheduleTimeoutEvent(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.clearTimeoutEvent", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ClearTimeoutEvent(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "runtime.getRandomData", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.GetRandomData(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.finalizeRef", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.FinalizeRef(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.stringVal", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.StringVal(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueGet", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueGet(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueSet", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueSet(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueDelete", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueDelete(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueIndex", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueIndex(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueSetIndex", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueSetIndex(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueCall", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueCall(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueInvoke", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueCall(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueNew", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueNew(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueLength", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueLength(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valuePrepareString", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValuePrepareString(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueLoadString", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueLoadString(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.valueInstanceOf", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.ValueInstanceOf(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.copyBytesToGo", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.CopyBytesToGo(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "syscall/js.copyBytesToJS", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.CopyBytesToJS(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	define("go", "debug", wasmtime.NewFunc(
		store,
		wasmtime.NewFuncType([]*wasmtime.ValType{i32}, []*wasmtime.ValType{}),
		func(caller *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
			mod.Debug(args[0].I32())
			return []wasmtime.Val{}, nil
		}),
	)

	return err
}
