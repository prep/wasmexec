package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var runtime = flag.String("runtime", "", "The host runtime")

func main() {
	flag.Parse()

	fmt.Printf("CMD: %s -runtime=%s %s\n", os.Args[0], *runtime, strings.Join(flag.Args(), " "))

	for _, env := range os.Environ() {
		fmt.Printf("ENV: %s\n", env)
	}
}
