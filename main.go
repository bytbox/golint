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
	for _, filename := range opts.Args {
		DoLint(filename)
	}
}

// Show version information
func ShowVersion() {
	fmt.Printf("golint v%s\n",version)
}

func DoLint(filename string) {

}

type StatelessLinter interface {
	Lint(string) (string, bool)
}

type StatefullLinter interface {
	Lint(string) (string, int, bool)
}

type LineLengthLint struct {}
func (lint LineLengthLint) Lint(line string) (msg string, err bool) {
 	length := 0
	// count characters - a tab is eight characters
	if length >= 80 {
		msg = fmt.Sprintf("line too long (%d >= %d)",length,80)
		err = true
	}
	return
}

