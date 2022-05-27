package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/prep/wasmexec"
	"github.com/prep/wasmexec/wasmerexec"

	"github.com/wasmerio/wasmer-go/wasmer"
)

var progname = filepath.Base(os.Args[0])

type Instance struct {
	wasmexec.Memory
	*wasmer.Instance

	spFn     wasmer.NativeFunction
	resumeFn wasmer.NativeFunction
}

func (instance *Instance) Debug(format string, params ...interface{}) {
}

func (instance *Instance) Error(format string, params ...interface{}) {
	log.Printf("ERROR: "+format+"\n", params...)
}

func (instance *Instance) GetSP() (int32, error) {
	val, err := instance.spFn()
	if err != nil {
		return 0, err
	}

	sp, ok := val.(int32)
	if !ok {
		return 0, fmt.Errorf("getsp: %T: expected an int32 return value", sp)
	}

	return sp, nil
}

func (instance *Instance) Resume() error {
	_, err := instance.resumeFn()
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

func run(filename string) error {
	// Read the Wasm file into memory.
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Create the engine and store.
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)
	defer store.Close()

	// Create a module out of the Wasm file.
	module, err := wasmer.NewModule(store, data)
	if err != nil {
		return err
	}
	defer module.Close()

	// Create an instance. This is needed here because the imports need an
	// instance to refer to, even though at this point no instance exists yet.
	instance := &Instance{}

	// Import the wasmexec functions and create an instance of the modules.
	imports := wasmerexec.Import(store, instance)
	instance.Instance, err = wasmer.NewInstance(module, imports)
	if err != nil {
		return err
	}

	// Fetch the memory export and set it on the instance, making the memory
	// accessible by the imports.
	mem, err := instance.Exports.GetMemory("mem")
	if err != nil {
		return err
	}

	instance.Memory = mem.Data()

	// Fetch the getsp function and reference it on the instance.
	instance.spFn, err = instance.Exports.GetFunction("getsp")
	if err != nil {
		return err
	}

	// Fetch the resume function and reference it on the instance.
	instance.resumeFn, err = instance.Exports.GetFunction("resume")
	if err != nil {
		return err
	}

	runFn, err := instance.Exports.GetFunction("run")
	if err != nil {
		return err
	}

	_, err = runFn(0, 0)
	return err
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
