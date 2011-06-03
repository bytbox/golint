// Copyright 2011 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

var version = "0.1.0"

var (
	verbose     = flag.Bool("v", false, "verbose output")
	showversion = flag.Bool("version", false, "show version information")
)

func main() {
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
	files = FilterSuffix(".go", ExpandFiles(args))
	if *verbose {
		for _, f := range files {
			fmt.Printf("%s ", f)
		}
		fmt.Printf("\n")
	}

	errs := make(chan os.Error)
	LintFiles(files, errs)
	close(errs)

	// Complain about any errors
	for err := range errs {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
}

func LintFiles(files []string, errs chan os.Error) {
	linelintDone := make(chan int)

	var err os.Error
	lintChan := make([]chan Lint, len(LineLinters))
	lineChan := make([]chan Line, len(LineLinters))
	lintRoot := make(chan Lint)
	for i, ll := range LineLinters {
		lineChan[i] = make(chan Line)
		lintChan[i] = make(chan Lint)
		go ll.RunLint(lineChan[i], lintChan[i])
	}

	// line-lint all files
	go func() {
		for _, fname := range files {
			var lines []string
			if lines, err = ReadFileLines(fname); err != nil {
				errs <- err
				continue
			}
			for lineno, line := range lines {
				for _, c := range lineChan {
					c <- Line{Location{fname, lineno}, line}
				}
			}
		}
		for _, c := range lineChan {
			close(c)
		}
	}()

	for _, c := range lintChan {
		go func() {
			for lint := range c {
				lintRoot <- lint
			}
		}()
	}

	go func() {
		for lint := range lintRoot {
			fmt.Printf("%s\n", lint)
		}
		linelintDone <- 1
	}()

	<-linelintDone
}

// Types of linters

type Linter interface {
	String() string
}

type LinterName struct {
	Category    string
	Name        string
	Description string
}

func (ln LinterName) String() string {
	return fmt.Sprintf("%s:%s: %s",ln.Category, ln.Name, ln.Description)
}

// Represents a line of code. Describes both the location of the code and the
// code itself.
type Line struct {
	Location
	line string
}

// Represents a line-based linter (which may or may not hold state).
//
// As a general rule, LineLinter should only be used for linters which do /not/
// need to hold state - those linters that must hold state would generally
// benefit from i.e. parsing capabilities.
type LineLinter interface {
	String() string
	RunLint(chan Line, chan Lint)
}

// A line-based linter using regular expressions
type RegexLinter struct {
	LinterName
	Regex string
}

func (rl RegexLinter) String() string {
	return fmt.Sprintf("%s (%s)", rl.LinterName.String(), rl.Regex)
}

func (rl RegexLinter) RunLint(text chan Line, lints chan Lint) {
	for line := range text {
		if matches, _ := regexp.Match(rl.Regex, []byte(line.line)); matches {
			lints <- LineLint{rl, line.Location, ""}
		}
	}
	close(lints)
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

type LineLint struct {
	linter Linter
	Location
	issue  string
}

func (lint LineLint) String() string {
	return fmt.Sprintf("%s at %s", lint.linter.String(), lint.Location)
}
