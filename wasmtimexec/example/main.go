package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/prep/wasmexec"
	"github.com/prep/wasmexec/wasmtimexec"

	"github.com/bytecodealliance/wasmtime-go"
)

var progname = filepath.Base(os.Args[0])

type Instance struct {
	wasmexec.Memory
	*wasmtime.Instance
	store *wasmtime.Store

	spFn     *wasmtime.Func
	resumeFn *wasmtime.Func
}

func (instance *Instance) Debug(format string, params ...interface{}) {
}

func (instance *Instance) Error(format string, params ...interface{}) {
	log.Printf("ERROR: "+format+"\n", params...)
}

func (instance *Instance) GetSP() (uint32, error) {
	val, err := instance.spFn.Call(instance.store)
	if err != nil {
		return 0, err
	}

	sp, ok := val.(int32)
	if !ok {
		return 0, fmt.Errorf("getsp: %T: expected an int32 return value", sp)
	}

	return uint32(sp), nil
}

func (instance *Instance) Resume() error {
	_, err := instance.resumeFn.Call(instance.store)
	return err
}

func (instance *Instance) Write(fd int, b []byte) (n int, err error) {
	switch fd {
	case 1:
		n, err = os.Stdout.Write(b)
	case 2:
		n, err = os.Stderr.Write(b)
	default:
		err = fmt.Errorf("%d: invalid file descriptor", fd)
	}

	return n, err
}

func (instance *Instance) Exit(code int) {
}

// HostCall is an optional method that allows the guest to use a fake waPC
// interface package to send messages to this host.
func (instance *Instance) HostCall(binding, namespace, operation string, payload []byte) ([]byte, error) {
	if namespace == "sample" && operation == "hello" {
		return []byte(fmt.Sprintf("Hello %s!", string(payload))), nil
	}

	return nil, fmt.Errorf("%s/%s: unsupported namespace + operation combination", namespace, operation)
}

func run(filename string) error {
	// Read the Wasm file into memory.
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Create the engine and store.
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	// Create a module out of the Wasm file.
	module, err := wasmtime.NewModule(engine, data)
	if err != nil {
		return err
	}

	// Create an instance. This is needed here because the imports need an
	// instance to refer to, even though at this point no instance exists yet.
	instance := &Instance{store: store}

	// Create the linker and import the wasmexec functions.
	linker := wasmtime.NewLinker(engine)
	gomod, err := wasmtimexec.Import(store, linker, instance)
	if err != nil {
		return err
	}

	// Create an instance of the module.
	if instance.Instance, err = linker.Instantiate(store, module); err != nil {
		return err
	}

	// Fetch the memory export and set it on the instance, making the memory
	// accessible by the imports.
	ext := instance.GetExport(store, "mem")
	if ext == nil {
		return errors.New("unable to find memory export")
	}

	mem := ext.Memory()
	if mem == nil {
		return errors.New("mem: export is not memory")
	}

	instance.Memory = wasmexec.NewMemory(mem.UnsafeData(store))

	// Fetch the getsp function and reference it on the instance.
	spFn := instance.GetExport(store, "getsp")
	if spFn == nil {
		return errors.New("getsp: missing export")
	}

	if instance.spFn = spFn.Func(); instance.spFn == nil {
		return errors.New("getsp: export is not a function")
	}

	// Fetch the resume function and reference it on the instance.
	resumeFn := instance.GetExport(store, "resume")
	if resumeFn == nil {
		return errors.New("resume: missing export")
	}

	if instance.resumeFn = resumeFn.Func(); instance.resumeFn == nil {
		return errors.New("resume: export is not a function")
	}

	// Fetch the "run" function and call it. This starts the program.
	runFn := instance.GetFunc(store, "run")
	if runFn == nil {
		return errors.New("run: missing export")
	}

	_, err = runFn.Call(store, 0, 0)
	if err != nil {
		return err
	}

	// Silently fail, because not all examples implement waPC.
	result, err := gomod.Invoke(nil, "hello", []byte("Host"))
	if err == nil {
		fmt.Printf("Message from guest: %s\n", string(result))
	}

	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [flags] program.wasm\n", progname)
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	if err := run(filename); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s: %s\n", progname, filename, err)
		os.Exit(1)
	}
}
