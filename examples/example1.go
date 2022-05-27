package main

import "os"

func main() {
	_, _ = os.Stdout.Write([]byte("Hello from Go!\n"))
}
