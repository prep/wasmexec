package main

import (
	"fmt"

	"github.com/prep/wasmexec/wapc"
)

func hello(payload []byte) ([]byte, error) {
	resp, err := wapc.HostCall("myBinding", "sample", "hello", []byte("Guest"))
	if err == nil {
		fmt.Printf("Message from host: %s\n", string(resp))
	}

	return []byte(fmt.Sprintf("Greetings %s!", string(payload))), nil
}

func main() {
	wapc.RegisterFunctions(wapc.Functions{
		"hello": hello,
	})
}
