// Copyright 2011 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

var version = "0.1.0"

var (
	verbose     = flag.Bool("v", false, "verbose output")
	showversion = flag.Bool("version", false, "show version information")
)

func main() {
	errs := make(chan os.Error)

	flag.Parse()
	args := flag.Args()

	if *showversion {
		fmt.Printf("golint %s\n", version)
		return
	}

	// find all files
	if *verbose {
		fmt.Printf("Scanning for source files... ")
	}
	files := args
	if len(files) == 0 {
		// just use the current directory if no files were specified
		files = make([]string, 1)
		files[0] = "."
	}
	files = FilterSuffix(".go", ExpandFiles(files))
	if *verbose {
		for _, f := range files {
			fmt.Printf("%s ", f)
		}
		fmt.Printf("\n")
	}

	go func() {
		// Complain about any errors
		for err := range errs {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}()

	LintFiles(files, errs)
	close(errs)

}

func LintFiles(files []string, errs chan os.Error) {
	lintWG := new(sync.WaitGroup)
	lintWG.Add(1)

	lintRoot := make(chan Lint)
	lintDone := make(chan int)

	var err os.Error

	// start up the LineLinters
	lineChan := make([]chan Line, len(LineLinters))
	for i, ll := range LineLinters {
		lineChan[i] = make(chan Line)
		go ll.RunLint(lineChan[i], lintRoot, lintWG)
	}

	go func() { // close the lintRoot when all linters report done
		lintWG.Wait()
		close(lintRoot)
	}()

	go func() { // Print out lints
		for lint := range lintRoot {
			fmt.Printf("%s\n", lint)
		}
		lintDone <- 1
	}()

	// lint all files
	for _, fname := range files {
		// line lint
		var lines []string
		if lines, err = ReadFileLines(fname); err != nil {
			errs <- err
			continue
		}
		for _, c := range lineChan {
			for lineno, line := range lines {
				c <- Line{Location{fname, lineno+1}, line}
			}
		}

		// parsing lint
		go func() {
			lintWG.Add(1)
			RunParsingLinters(fname, lintRoot, errs)
			lintWG.Done()
		}()
	}

	for _, c := range lineChan { // close all lineChans
		close(c)
	}

	// we're done putting data into the linters
	lintWG.Done()
	<-lintDone // wait for printing to cease

}

// Types of linters

type Linter interface {
	String() string
}

// Name and description of a linter.
type LinterName struct {
	Category    string
	Name        string
	Description string
}

func (ln LinterName) String() string {
	return fmt.Sprintf("%s:%s: %s", ln.Category, ln.Name, ln.Description)
}

// Represents a line of code. Describes both the location of the code and the
// code itself.
type Line struct {
	Location
	line string
}

type Lint interface {
	String() string
}

// Represents a location in a file
type Location struct {
	filename string
	lineno int
}

func (loc Location) String() string {
	return fmt.Sprintf("%s:%d", loc.filename, loc.lineno)
}

