package main

import (
	"fmt"
	"opts"
	"os"
)

var version = "0.0.1"

// options
var showVersion = opts.Longflag("version",
	"display version information and exit") 

func main() {
	// Do the argument parsing
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}

}

// Show version information
func ShowVersion() {
	fmt.Printf("golint v%s\n",version)
}

