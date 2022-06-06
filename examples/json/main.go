package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"os"
)

type Message struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Photo   []byte `json:"photo"`
}

func main() {
	photo := make([]byte, 32)
	if _, err := rand.Read(photo); err != nil {
		log.Fatalf("Unable to generate random data: %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	_ = enc.Encode(Message{
		Title:   "Extra! Extra!",
		Message: "Go Wasm binaries now run in wasmer, wasmtime and wazero!",
		Photo:   photo,
	})
}
