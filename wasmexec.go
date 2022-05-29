package wasmexec

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"reflect"
	"syscall"
	"time"
)

// enosys error code meaning "function not supported".
const enosys = uint32(52)

// nanHead is the NaN-header for values that are not a number, but an ID.
const nanHead = 0x7FF80000

// NaN describes a not-a-number value.
var NaN = math.NaN()

// Properties describe the properties on an object.
type Properties map[string]any

// Object describes a JSON object.
type Object struct {
	properties Properties
}

// Array describes an array of elements.
type Array struct {
	elements []any
}

// Uint8Array describes a byte slice.
type Uint8Array struct {
	data []byte
}

// String represents a stored string.
type String struct {
	payload string
}

// Function describes a function.
type Function struct {
	fn func(args []any) any
}

// newFuncObject returns a new function.
func newFuncObject(fn func(args []any) any) *Function {
	return &Function{fn: fn}
}

func toInt(v any) (int, error) {
	if val, ok := v.(int); ok {
		return val, nil
	}

	if val, ok := v.(float64); ok {
		return int(val), nil
	}

	return 0, fmt.Errorf("%T: unable to convert to int", v)
}

// ModuleGo implements the JavaScript imports that a Go program compiled with
// GOOS=js expects.
type ModuleGo struct {
	instance  Instance
	idcounter uint32
	ids       map[any]uint32
	values    map[uint32]any
	refcounts map[uint32]int32
}

// NewModuleGo returns a new ModuleGo.
func NewModuleGo(instance Instance) *ModuleGo {
	var mod *ModuleGo

	mod = &ModuleGo{
		instance:  instance,
		idcounter: 10,
		refcounts: make(map[uint32]int32),
		ids: map[any]uint32{
			NaN:        0,
			float64(0): 1,
			nil:        2,
			true:       3,
			false:      4,
		},
		values: map[uint32]any{
			0: NaN,
			1: 0,
			2: nil,
			3: true,
			4: false,

			// global.
			5: &Object{
				properties: Properties{
					"Array": &Function{
						fn: func([]any) any {
							return &Array{}
						},
					},

					"Date": &Function{
						fn: func([]any) any {
							return &Object{
								properties: Properties{
									"getTimezoneOffset": &Function{
										fn: func(args []any) any {
											t := time.Now()
											_, offset := t.Zone()
											return (offset / 60) * -1
										},
									},
								},
							}
						},
					},

					"Object": &Function{
						fn: func([]any) any {
							return &Object{properties: make(Properties)}
						},
					},

					"Uint8Array": &Function{
						fn: func(args []any) any {
							if len(args) == 0 {
								return []byte{}
							}

							length, err := toInt(args[0])
							if err != nil {
								return []byte{}
							}

							return &Uint8Array{
								data: make([]byte, length),
							}
						},
					},

					"crypto": &Object{
						properties: Properties{
							"getRandomValues": &Function{
								fn: func(args []any) any {
									if len(args) != 1 {
										mod.instance.Error("crypto.getRandomValues: %d: invalid number of arguments", len(args))
										return 0
									}

									a, ok := args[0].(*Uint8Array)
									if !ok {
										mod.instance.Error("crypto.getRandomValues: %T: not type Uint8Array", args[0])
										return 0
									}

									n, err := rand.Read(a.data)
									if err != nil {
										mod.instance.Error("crypto.getRandomValues: %v", err)
										return 0
									}

									return n
								},
							},
						},
					},

					"fs": &Object{
						properties: Properties{
							"constants": Properties{
								"O_WRONLY": syscall.O_WRONLY,
								"O_RDWR":   syscall.O_RDWR,
								"O_CREAT":  syscall.O_CREAT,
								"O_TRUNC":  syscall.O_TRUNC,
								"O_APPEND": syscall.O_APPEND,
								"O_EXCL":   syscall.O_EXCL,
							},

							"write": &Function{
								fn: func(args []any) any {
									if len(args) != 6 {
										mod.instance.Error("fs.write: %d: invalid number of arguments", len(args))
										return nil
									}

									val, ok := args[0].(float64)
									if !ok {
										mod.instance.Error("fs.write: %T: not type float64", args[0])
										return nil
									}
									fd := int(val)

									buf, ok := args[1].(*Uint8Array)
									if !ok {
										mod.instance.Error("fs.write: %T: not type Uint8Array", args[1])
										return nil
									}

									/*
										offset, ok := args[2].(int)
										if !ok {
											mod.instance.Error("fs.write: %T: not type int", args[2])
											return 9
										}

										val, ok = args[3].(float64)
										if !ok {
											mod.instance.Error("fs.write: %T: not type float64", args[3])
											return 9
										}
										length := int(val)

										var position int64
										if args[4] != nil {
											val, ok = args[4].(float64)
											if !ok {
												mod.instance.Error("fs.write: %T: not type float64", args[4])
												return 9
											}

											position = int64(val)
										}

										_ = offset
										_ = length
										_ = position
									*/

									callback, ok := args[5].(*Function)
									if !ok {
										mod.instance.Error("fs.write: %T: not type Function", args[5])
										return nil
									}

									n, err := mod.instance.Write(fd, buf.data)
									if err != nil {
										// TODO: this does not work.
										callback.fn([]any{enosys})
										return nil
									}

									callback.fn([]any{nil, n})
									return nil
								},
							},
						},
					},

					"process": &Object{
						properties: Properties{
							"getuid":    newFuncObject(func([]any) any { return -1 }),
							"getgid":    newFuncObject(func([]any) any { return -1 }),
							"geteuid":   newFuncObject(func([]any) any { return -1 }),
							"getegid":   newFuncObject(func([]any) any { return -1 }),
							"getgroups": newFuncObject(func([]any) any { return enosys }),
							"pid":       -1,
							"ppid":      -1,
							"umask":     newFuncObject(func([]any) any { return enosys }),
							"cwd":       newFuncObject(func([]any) any { return enosys }),
							"chdir":     newFuncObject(func([]any) any { return enosys }),
						},
					},
				},
			},

			// jsGo.
			6: &Object{
				properties: Properties{
					"_pendingEvent": nil,

					// This is called by js.FuncOf().
					"_makeFuncWrapper": &Function{
						fn: func(args []any) any {
							if len(args) == 0 {
								return nil
							}

							id := args[0]

							return &Function{
								fn: func(args []any) any {
									event := &Object{
										properties: Properties{
											"id":   id,
											"this": mod.values[6].(*Object),
											"args": &Array{elements: args},
										},
									}

									mod.values[6].(*Object).properties["_pendingEvent"] = event
									if err := mod.instance.Resume(); err != nil {
										mod.instance.Error("_makeFuncWrapper: Resume: %v", err)
										return nil
									}

									return event.properties["result"]
								},
							}
						},
					},
				},
			},
		},
	}

	return mod
}

// ****************************************************************************
// **************************** [ Helper methods ] ****************************
// ****************************************************************************

// TODO: Perhaps we need a better scheme of assigning IDs to in-memory objects.
func (mod *ModuleGo) getID() uint32 {
	id := mod.idcounter
	mod.idcounter++

	return id
}

// getInt32 returns an int32 value.
func (mod *ModuleGo) getInt32(offset int32) (int32, error) {
	val, err := mod.instance.GetUInt32(offset)
	if err != nil {
		return 0, err
	}

	return int32(val), nil
}

// setInt32 sets an int32 value.
func (mod *ModuleGo) setInt32(offset, val int32) error {
	return mod.instance.SetUInt32(offset, uint32(val))
}

// loadValue loads either a number from the specified address, or it loads an
// object ID from the address and fetches that value from the stored values.
func (mod *ModuleGo) loadValue(addr int32) (any, error) {
	f, err := mod.instance.GetFloat64(addr)
	switch {
	case err != nil:
		return nil, err
	case f == 0:
		return nil, nil
	case !math.IsNaN(f):
		return f, nil
	}

	id, err := mod.instance.GetUInt32(addr)
	if err != nil {
		return nil, err
	}

	mod.instance.Debug("   loadValue(id=%v)", id)

	return mod.values[id], nil
}

func (mod *ModuleGo) storeValue(addr int32, v any) error {
	mod.instance.Debug("   storeValue(addr=%v v=%v nil=%v)", addr, v, (v == nil))

	// Convert any integer to a float64, which is akin to a JSON number.
	switch num := v.(type) {
	case int:
		v = float64(num)
	case uint:
		v = float64(num)
	case int8:
		v = float64(num)
	case uint8:
		v = float64(num)
	case int16:
		v = float64(num)
	case uint16:
		v = float64(num)
	case int32:
		v = float64(num)
	case uint32:
		v = float64(num)
	case int64:
		v = float64(num)
	case uint64:
		v = float64(num)
	case float32:
		v = float64(num)
	}

	// setNaN sets a NaN-value on the specified address.
	setNaN := func(val uint32) error {
		if err := mod.instance.SetUInt32(addr+4, nanHead); err != nil {
			return err
		}

		return mod.instance.SetUInt32(addr, val)
	}

	// If this is a number, store it as such.
	if num, ok := v.(float64); ok && num != 0 {
		if math.IsNaN(num) {
			return setNaN(0)
		}

		return mod.instance.SetFloat64(addr, num)
	}

	// Check for specific values that don't require storing anything in the
	// ids and values map.
	switch v {
	case nil:
		return setNaN(2)
	case true:
		return setNaN(3)
	case false:
		return setNaN(4)
	}

	// Create a unique signature of the value.
	signature := fmt.Sprintf("%d", reflect.ValueOf(v).Pointer())
	mod.instance.Debug("   storeValue(type=%T signature=%v)", v, signature)

	// Use the signature to check if this value has already been stored. If not,
	// store it in the ids and values map.
	id, ok := mod.ids[signature]
	if !ok {
		id = mod.getID()
		mod.values[id] = v
		mod.refcounts[id] = 0
		mod.ids[signature] = id
	}

	// Raise the reference count on this object.
	mod.refcounts[id]++

	// Determine if the value needs to be stored with a specific type flag.
	var typeFlag uint32
	switch t := v.(type) {
	case *Object, *Array, *Uint8Array, Properties:
		if t != nil {
			typeFlag = 1
		}

	case *String:
		typeFlag = 2

	case *Function:
		typeFlag = 4
	default:
		panic(fmt.Sprintf("%T: unknown value type", t))
	}

	mod.instance.Debug("   storeValue(id=%v typeFlag=%v refcount=%v signature=%q)", id, typeFlag, mod.refcounts[id], signature)

	// Store the type.
	if err := mod.instance.SetUInt32(addr+4, nanHead|typeFlag); err != nil {
		return err
	}

	// Store the ID.
	return mod.instance.SetUInt32(addr, id)
}

// loadSlice returns a byte slice that is referenced by the specified address.
func (mod *ModuleGo) loadSlice(addr int32) ([]byte, error) {
	a, err := mod.instance.GetInt64(addr)
	if err != nil {
		return nil, err
	}

	length, err := mod.instance.GetInt64(addr + 8)
	if err != nil {
		return nil, err
	}

	mod.instance.Debug("   loadSlice(array=%v length=%v)", a, length)

	return mod.instance.Mem(int32(a), int32(length))
}

// loadSliceOfValues returns a slice of values that is referenced by the
// specified address.
func (mod *ModuleGo) loadSliceOfValues(addr int32) ([]any, error) {
	offset, err := mod.instance.GetInt64(addr)
	if err != nil {
		return nil, err
	}

	length, err := mod.instance.GetInt64(addr + 8)
	if err != nil {
		return nil, err
	}

	a := make([]any, length)
	for i := int64(0); i < length; i++ {
		a[i], err = mod.loadValue(int32(offset + (i * 8)))
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

// loadString returns a string that is referenced by the specified address.
func (mod *ModuleGo) loadString(addr int32) (string, error) {
	d, err := mod.loadSlice(addr)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

func (mod *ModuleGo) reflectApply(v any, name string, args []any) (any, error) {
	mod.instance.Debug("   reflectApply(name=%v)", name)

	obj, err := mod.reflectGet(v, name)
	if err != nil {
		return nil, err
	}

	return mod.reflectConstruct(obj, args)
}

func (mod *ModuleGo) reflectConstruct(v any, args []any) (any, error) {
	mod.instance.Debug("   reflectConstruct(v=%v args=%v)", v, args)

	if fn, ok := v.(*Function); ok {
		return fn.fn(args), nil
	}

	return nil, fmt.Errorf("%T: not a function", v)
}

func (mod *ModuleGo) reflectGet(v, key any) (any, error) {
	mod.instance.Debug("   reflectGet(key=%v)", key)

	if v == nil {
		v = mod.values[5]
	}

	if name, ok := key.(string); ok {
		switch vv := v.(type) {
		case *Object:
			return vv.properties[name], nil
		case Properties:
			return vv[name], nil
		}
	}

	index, ok := key.(int64)
	if !ok {
		return nil, errors.New("key not an int64")
	}

	a, ok := v.(*Array)
	switch {
	case !ok:
		return nil, errors.New("value not a slice")
	case index < 0 || index >= int64(len(a.elements)):
		return nil, errors.New("index out of range")
	}

	return a.elements[index], nil
}

func (mod *ModuleGo) reflectSet(v, key, value any) error {
	mod.instance.Debug("   reflectSet(v=%v key=%v value=%v)", v, key, value)

	if v == nil {
		v = mod.values[5]
	}

	if name, ok := key.(string); ok {
		v.(*Object).properties[name] = value
		return nil
	}

	index, ok := key.(int64)
	if !ok {
		return errors.New("key not an int64")
	}

	a, ok := v.(*Array)
	switch {
	case !ok:
		return errors.New("value not a slice")
	case index < 0 || index >= int64(len(a.elements)):
		return errors.New("index out of range")
	}

	a.elements[index] = value
	return nil
}

func (mod *ModuleGo) reflectDeleteProperty(v, key any) error {
	mod.instance.Debug("   reflectDelete(v=%v key=%v)", v, key)

	if v == nil {
		v = mod.values[5]
	}

	if name, ok := key.(string); ok {
		delete(v.(*Object).properties, name)
		return nil
	}

	index, ok := key.(int64)
	if !ok {
		return errors.New("key not an int64")
	}

	a, ok := v.(*Array)
	switch {
	case !ok:
		return errors.New("value not a slice")
	case index < 0 || index >= int64(len(a.elements)):
		return errors.New("index out of range")
	}

	copy(a.elements[index:], a.elements[index+1:])
	a.elements[len(a.elements)-1] = nil
	a.elements = a.elements[:len(a.elements)-1]

	return nil
}

func (mod *ModuleGo) wrap(name string, fn func() error) error {
	if fn == nil {
		mod.instance.Error("%s NOT IMPLEMENTED", name)
		return nil
	}

	if name != "" {
		mod.instance.Debug(name)
	}

	if err := fn(); err != nil {
		if name != "" {
			mod.instance.Error("%s: %v", name, err)
		}
		return err
	}

	return nil
}

// ****************************************************************************
// ***************************** [ Go JS module ] *****************************
// ****************************************************************************

// WasmExit is called whenever the WASM program exits.
//
// This method is called from the runtime package.
func (mod *ModuleGo) WasmExit(sp int32) {
	sp >>= 0

	_ = mod.wrap("runtime.wasmExit", func() error {
		v, err := mod.instance.GetUInt32(sp + 8)
		if err != nil {
			return err
		}

		mod.instance.Exit(int(v))
		return nil
	})
}

// WasmWrite writes data to a file descriptor.
//
// This method is called from the runtime package.
func (mod *ModuleGo) WasmWrite(sp int32) {
	sp >>= 0

	_ = mod.wrap("runtime.wasmWrite", func() error {
		fd, err := mod.instance.GetInt64(sp + 8)
		if err != nil {
			return err
		}

		p, err := mod.instance.GetInt64(sp + 16)
		if err != nil {
			return err
		}

		n, err := mod.getInt32(sp + 24)
		if err != nil {
			return err
		}

		mem, err := mod.instance.Mem(int32(p), n)
		if err != nil {
			return err
		}

		_, err = mod.instance.Write(int(fd), mem)
		return err
	})
}

// ResetMemoryDataView is called whenever WebAssembly's memory.grow instruction
// has been used.
//
// This method is called from the runtime package.
func (mod *ModuleGo) ResetMemoryDataView(sp int32) {
	sp >>= 0

	_ = mod.wrap("runtime.resetMemoryDataView", nil)
}

// Nanotime1 returns the current time in nanoseconds.
//
// This method is called from the runtime package.
func (mod *ModuleGo) Nanotime1(sp int32) {
	sp >>= 0

	_ = mod.wrap("runtime.nanotime1", func() error {
		return mod.instance.SetInt64(sp+8, time.Now().UnixNano())
	})
}

// Walltime returns the current seconds and nanoseconds.
//
// This method is called from the runtime package.
func (mod *ModuleGo) Walltime(sp int32) {
	sp >>= 0

	_ = mod.wrap("runtime.walltime", func() error {
		msec := time.Now().UnixNano() / int64(time.Millisecond)

		if err := mod.instance.SetInt64(sp+8, msec/1000); err != nil {
			return err
		}

		return mod.setInt32(sp+16, int32(msec%1000)*1000000)
	})
}

// ScheduleTimeoutEvent is called whenever an event needs to be scheduled after
// a certain amount of milliseconds.
//
// This method is called from the runtime package.
func (mod *ModuleGo) ScheduleTimeoutEvent(sp int32) {
	sp >>= 0

	_ = mod.wrap("runtime.scheduleTimeoutEvent", nil)
}

// ClearTimeoutEvent clears a timeout event scheduled by ScheduleTimeoutEvent.
//
// This method is called from the runtime package.
func (mod *ModuleGo) ClearTimeoutEvent(sp int32) {
	sp >>= 0

	_ = mod.wrap("runtime.clearTimeoutEvent", nil)
}

// GetRandomData returns random data.
//
// This method is called from the runtime package.
func (mod *ModuleGo) GetRandomData(sp int32) {
	sp >>= 0

	_ = mod.wrap("runtime.getRandomData", func() error {
		data, err := mod.loadSlice(sp + 8)
		if err != nil {
			return err
		}

		_, err = rand.Read(data)
		return err
	})
}

// FinalizeRef removes a value from memory.
//
// This method is called from various places in syscall/js.Value.
func (mod *ModuleGo) FinalizeRef(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.finalizeRef", func() error {
		id, err := mod.instance.GetUInt32(sp + 8)
		if err != nil {
			return err
		}

		// Make sure the ID has a reference count.
		ref, ok := mod.refcounts[id]
		if !ok {
			return fmt.Errorf("%d: missing reference count for id", id)
		}

		ref--

		// If the reference count is 0, clean up the object.
		if ref == 0 {
			signature, ok := mod.values[id]
			if !ok {
				return fmt.Errorf("%d: could not find signature in values for id", id)
			}

			mod.instance.Debug("%d: deleting object", id)

			delete(mod.refcounts, id)
			delete(mod.values, id)
			delete(mod.ids, signature)
		}

		return nil
	})
}

// StringVal stores a value as a string.
//
// This method is called from syscall/js.ValueOf().
func (mod *ModuleGo) StringVal(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.stringVal", func() error {
		v, err := mod.loadString(sp + 8)
		if err != nil {
			return err
		}

		return mod.storeValue(sp+24, &String{v})
	})
}

// ValueGet returns the JavaScript property of an object.
//
// This method is called from syscall/js.Value.Get().
func (mod *ModuleGo) ValueGet(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.valueGet", func() error {
		// Fetch the object.
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		// Fetch the name of the property to read.
		name, err := mod.loadString(sp + 16)
		if err != nil {
			return err
		}

		// Read the value from the property on the object.
		result, err := mod.reflectGet(v, name)
		if err != nil {
			return err
		}

		resultSP, err := mod.instance.GetSP()
		if err != nil {
			return err
		}

		resultSP >>= 0

		// Store the results.
		return mod.storeValue(resultSP+32, result)
	})
}

// ValueSet sets a value on a property on an object.
//
// This method is called from syscall/js.Value.Set().
func (mod *ModuleGo) ValueSet(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.valueSet", func() error {
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		name, err := mod.loadString(sp + 16)
		if err != nil {
			return err
		}

		value, err := mod.loadValue(sp + 32)
		if err != nil {
			return err
		}

		return mod.reflectSet(v, name, value)
	})
}

// ValueDelete deletes a property on an object.
//
// This method is called from syscall/js.Value.Delete().
func (mod *ModuleGo) ValueDelete(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.valueDelete", func() error {
		// Fetch the object.
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		// Fetch the property name.
		name, err := mod.loadString(sp + 16)
		if err != nil {
			return err
		}

		// Delete the property on the object.
		return mod.reflectDeleteProperty(v, name)
	})
}

// ValueIndex returns a value at a particular index in an array.
//
// This method is called from syscall/js.Value.Index().
func (mod *ModuleGo) ValueIndex(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.valueIndex", func() error {
		// Fetch the object.
		obj, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		// Fetch the index to fetch the value for.
		index, err := mod.instance.GetInt64(sp + 16)
		if err != nil {
			return err
		}

		// Fetch the value on the index in the array.
		result, err := mod.reflectGet(obj, index)
		if err != nil {
			return err
		}

		// Return the value.
		return mod.storeValue(sp+24, result)
	})
}

// ValueSetIndex sets the value at a particular index of an array.
//
// This method is called from syscall/js.Value.SetIndex().
func (mod *ModuleGo) ValueSetIndex(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.valueSetIndex", func() error {
		// Fetch the object.
		obj, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		// Fetch the index.
		index, err := mod.instance.GetInt64(sp + 16)
		if err != nil {
			return err
		}

		// Fetch the value that should be set on the index in the array.
		value, err := mod.loadValue(sp + 24)
		if err != nil {
			return err
		}

		// Set the value on the index in the array.
		return mod.reflectSet(obj, index, value)
	})
}

// ValueCall calls the method on an object with the give arguments.
//
// This method is called from syscall/js.Value.Call().
func (mod *ModuleGo) ValueCall(sp int32) {
	sp >>= 0

	var resultSP int32
	err := mod.wrap("syscall/js.valueCall", func() error {
		var err error
		resultSP, err = mod.instance.GetSP()
		if err != nil {
			return err
		}

		resultSP >>= 0

		// Fetch the object.
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		// Fetch the name of the method to call.
		name, err := mod.loadString(sp + 16)
		if err != nil {
			return err
		}

		// Fetch the arguments to call the method with.
		args, err := mod.loadSliceOfValues(sp + 32)
		if err != nil {
			return err
		}

		// Call the method on the object with the arguments.
		result, err := mod.reflectApply(v, name, args)
		if err != nil {
			return err
		}

		// Store the results of the call.
		if err = mod.storeValue(resultSP+56, result); err != nil {
			return err
		}

		return mod.instance.SetUInt8(resultSP+64, 1)
	})
	if err == nil || resultSP == 0 {
		return
	}

	if err = mod.storeValue(resultSP+56, err); err != nil {
		return
	}

	_ = mod.instance.SetUInt8(resultSP+54, 0)
}

// ValueInvoke calls the value v with the specified arguments.
//
// This method is called from syscall/js.Value.Invoke().
func (mod *ModuleGo) ValueInvoke(sp int32) {
	sp >>= 0

	var resultSP int32
	err := mod.wrap("syscall/js.valueInvoke", func() error {
		var err error
		resultSP, err = mod.instance.GetSP()
		if err != nil {
			return err
		}

		resultSP >>= 0

		// Fetch the function.
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		// Fetch the arguments.
		args, err := mod.loadSliceOfValues(sp + 16)
		if err != nil {
			return err
		}

		// Call v with the specified arguments.
		result, err := mod.reflectConstruct(v, args)
		if err != nil {
			return err
		}

		// Store the results of the call.
		if err = mod.storeValue(resultSP+40, result); err != nil {
			return err
		}

		return mod.instance.SetUInt8(resultSP+48, 1)
	})
	if err == nil || resultSP == 0 {
		return
	}

	if err = mod.storeValue(resultSP+40, err); err != nil {
		return
	}

	_ = mod.instance.SetUInt8(resultSP+48, 0)
}

// ValueNew calls a constructor function with the given arguments. This is akin
// to JavaScript's "new" operator.
//
// This method is called from syscall/js.Value.New().
func (mod *ModuleGo) ValueNew(sp int32) {
	sp >>= 0

	var resultSP int32
	err := mod.wrap("syscall/js.valueNew", func() error {
		var err error
		resultSP, err = mod.instance.GetSP()
		if err != nil {
			return err
		}

		resultSP >>= 0

		// Fetch the constructor function.
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		// Fetch the arguments to call the constructor with.
		args, err := mod.loadSliceOfValues(sp + 16)
		if err != nil {
			return err
		}

		// Call the constructor function with the arguments.
		result, err := mod.reflectConstruct(v, args)
		if err != nil {
			return err
		}

		// Store the results of the call.
		if err = mod.storeValue(resultSP+40, result); err != nil {
			return err
		}

		return mod.instance.SetUInt8(resultSP+48, 1)
	})
	if err == nil || resultSP == 0 {
		return
	}

	if err = mod.storeValue(resultSP+40, err); err != nil {
		return
	}

	_ = mod.instance.SetUInt8(resultSP+48, 1)
}

// ValueLength returns the JavaScript property of "length" of v.
//
// This method is called from syscall/js.Value.Length().
func (mod *ModuleGo) ValueLength(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.valueLength", func() error {
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		// TODO: Supposedly this should only be called on objects?
		switch val := v.(type) {
		case *Array:
			return mod.instance.SetInt64(sp+16, int64(len(val.elements)))
		default:
			return fmt.Errorf("%T: unknown type for valueLength", v)
		}
	})
}

// ValuePrepareString converts a value to a string and stores it.
//
// This method is called from syscall/js.Value.String().
func (mod *ModuleGo) ValuePrepareString(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.valuePrepareString", func() error {
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		var s *String
		switch vv := v.(type) {
		case *Uint8Array:
			s = &String{payload: string(vv.data)}
		default:
			return fmt.Errorf("%T: unable to convert type to string", v)
		}

		if err = mod.storeValue(sp+16, s); err != nil {
			return err
		}

		return mod.instance.SetInt64(sp+24, int64(len(s.payload)))
	})
}

// ValueLoadString loads a string into memory.
//
// This method is called from syscall/js.Value.String().
func (mod *ModuleGo) ValueLoadString(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.valueLoadString", func() error {
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		s, ok := v.(*String)
		if !ok {
			return fmt.Errorf("%T: type not a string", v)
		}

		dst, err := mod.loadSlice(sp + 16)
		if err != nil {
			return err
		}

		copy(dst, s.payload)
		return nil
	})
}

// ValueInstanceOf returns true when v is an instance of type t.
//
// This method is called from syscall/js.Value.InstanceOf().
func (mod *ModuleGo) ValueInstanceOf(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.valueInstanceOf", func() error {
		/*
			v, err := mod.loadValue(sp + 8)
			if err != nil {
				return err
			}

			t, err := mod.loadValue(sp + 16)
			if err != nil {
				return err
			}

			mod.instance.Debug("   is %#v an instance of %#v?", src, dst)
		*/

		// TODO: This is not an implementation. Setting "0" means this method will always return false.
		return mod.instance.SetUInt8(sp+24, 0)
	})
}

// CopyBytesToGo copies bytes from JavaScript to Go.
func (mod *ModuleGo) CopyBytesToGo(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.copyBytesToGo", func() error {
		dst, err := mod.loadSlice(sp + 8)
		if err != nil {
			return err
		}

		v, err := mod.loadValue(sp + 32)
		if err != nil {
			return err
		}

		src, ok := v.(*Uint8Array)
		if !ok {
			return fmt.Errorf("src: %T not type Uint8Array", v)
		}

		if len(dst) == 0 || len(src.data) == 0 {
			return mod.instance.SetUInt8(sp+48, 0)
		}

		n := copy(dst, src.data)
		if err = mod.instance.SetInt64(sp+40, int64(n)); err != nil {
			return err
		}

		return mod.instance.SetUInt8(sp+48, 1)
	})
}

// CopyBytesToJS copies bytes from Go to JavaScript.
func (mod *ModuleGo) CopyBytesToJS(sp int32) {
	sp >>= 0

	_ = mod.wrap("syscall/js.copyBytesToJS", func() error {
		v, err := mod.loadValue(sp + 8)
		if err != nil {
			return err
		}

		dst, ok := v.(*Uint8Array)
		if !ok {
			return fmt.Errorf("dst: %T not type Uint8Array", v)
		}

		src, err := mod.loadSlice(sp + 16)
		if err != nil {
			return err
		}

		if len(dst.data) == 0 || len(src) == 0 {
			return mod.instance.SetUInt8(sp+48, 0)
		}

		n := copy(dst.data, src)
		if err = mod.instance.SetInt64(sp+40, int64(n)); err != nil {
			return err
		}

		return mod.instance.SetUInt8(sp+48, 1)
	})
}

// Debug prints some debugging information ... I guess?
func (mod *ModuleGo) Debug(sp int32) {
	sp >>= 0

	_ = mod.wrap("debug", nil)
}
