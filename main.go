package main

import (
	"fmt"
	"io/ioutil"
	"opts"
	"os"
	"strings"
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
		err := DoLint(filename)
		if err != nil {

		}
	}
}

// Show version information
func ShowVersion() {
	fmt.Printf("golint v%s\n",version)
}

var statelessLinters = []StatelessLinter {
	LineLengthLint{},
}

var statefullLinters = []StatefullLinter {

}

type StatelessLinter interface {
	Lint(string) (string, bool)
}

type StatefullLinter interface {
	Lint(string, int) (string, bool)
}

func DoLint(filename string) os.Error {
	// read in the file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// for each line
	lines := strings.Split(string(content), "\n", -1)
	for lineno, line := range lines {
		// run through the stateless linters
		for _, linter := range statelessLinters {
			msg, err := linter.Lint(line)
			if err {
				fmt.Printf("L%d: %s\n", lineno, msg)
			}
		}
		// run through the statefull linters
		for _, linter := range statefullLinters {
			msg, err := linter.Lint(line, lineno)
			if err {
				fmt.Printf("L%d: %s\n", lineno, msg)
			}
		}
	}
	return nil
}

