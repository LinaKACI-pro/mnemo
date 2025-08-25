package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "0.0.1"

func main() {
	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Println("mnemo version", version)
		os.Exit(0)
	}

	fmt.Println("Hello from mnemo!")
}
