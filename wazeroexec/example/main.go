package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/prep/wasmexec"
	"github.com/prep/wasmexec/wazeroexec"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

var progname = filepath.Base(os.Args[0])

type Instance struct {
	wasmexec.Memory
	spFn     api.Function
	resumeFn api.Function
}

func (instance *Instance) Debug(format string, params ...interface{}) {
	// log.Printf("DEBUG: "+format+"\n", params...)
}

func (instance *Instance) Error(format string, params ...interface{}) {
	log.Printf("ERROR: "+format+"\n", params...)
}

func (instance *Instance) GetSP() (uint32, error) {
	results, err := instance.spFn.Call(context.Background())
	switch {
	case err != nil:
		return 0, err
	case len(results) == 0:
		return 0, errors.New("getsp: no sp value returned")
	}

	return uint32(results[0]), nil
}

func (instance *Instance) Resume() error {
	_, err := instance.resumeFn.Call(context.Background())
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

	ctx := context.Background()

	// Create the runtime.
	runtime := wazero.NewRuntime()
	defer runtime.Close(ctx)

	// Create an instance. This is needed here because the imports need an
	// instance to refer to, even though at this point no instance exists yet.
	instance := &Instance{}

	// Import the wasmexec functions.
	gomod, err := wazeroexec.Import(ctx, runtime, runtime, instance)
	if err != nil {
		return err
	}

	// Create an instance of the module.
	module, err := runtime.InstantiateModuleFromBinary(ctx, data)
	if err != nil {
		return err
	}
	defer module.Close(ctx)

	// Fetch the memory export and set it on the instance, making the memory
	// accessible by the imports.
	mem := module.ExportedMemory("mem")
	if mem == nil {
		return errors.New("unable to find memory export")
	}

	instance.Memory = wazeroexec.NewMemory(mem)

	// Fetch the getsp function and reference it on the instance.
	instance.spFn = module.ExportedFunction("getsp")
	if instance.spFn == nil {
		return errors.New("getsp: missing export")
	}

	// Fetch the resume function and reference it on the instance.
	instance.resumeFn = module.ExportedFunction("resume")
	if instance.resumeFn == nil {
		return errors.New("resume: missing export")
	}

	// Fetch the "run" function and call it. This starts the program.
	runFn := module.ExportedFunction("run")
	if runFn == nil {
		return errors.New("run: missing export")
	}

	_, err = runFn.Call(ctx, 0, 0)
	if err != nil {
		return err
	}

	// Silently fail, because not all examples implement waPC.
	result, err := gomod.Invoke(ctx, "hello", []byte("Host"))
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
