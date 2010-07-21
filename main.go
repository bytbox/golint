/*
	package comment
*/
package main

import (
	"fmt"
	"flag"
	"os"
)

var version = "0.0.1"

var opts struct {
	showVersion *bool
	verbose *bool
}

func main() {
	// Option definitions
	opts.showVersion = flag.Bool("version", 
		false,
		"display version information and exit")
	opts.verbose = flag.Bool("verbose",
		false,
		"use verbose output")
	// Do the argument parsing
	flag.Parse()
	if *opts.showVersion {
		ShowVersion()
		os.Exit(0)
	}

}

// Show version information
func ShowVersion() {
	fmt.Printf("golint v%s\n",version)
}

// Log output in an appropriate fashion, depending on settings
func Message(msg string) {
	fmt.Printf("%s\n",msg)
}

// Log verbose output the right way, depending on settings
func Verbose(msg string) {
	if *opts.verbose {
		fmt.Printf("%s\n",msg)
	}
}
