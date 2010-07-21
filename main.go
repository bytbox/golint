/*
	package comment
*/
package main

import (
	"fmt"
	"flag"
	"os"
)

var version = '0.0.1'

var opts struct {
	showVersion *bool
}

func main() {
	// Option definitions
	opts.showVersion = flag.Bool("version", 
			false, 
			"display version information and exit")

	// Do the argument parsing
	flag.Parse()
	if *opts.showVersion {
		ShowVersion()
		os.Exit(0)
	}

	fmt.Printf("see!")
}

func ShowVersion() {
	fmt.Printf("golint v%s\n",version)
}
