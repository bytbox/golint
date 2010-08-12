package main

import (
	"regexp"
)

// TodoLint is a stateless linter that checks for and prints out lines that
// have "TODO" notices.
type TodoLint struct{}

func (TodoLint) Lint(line string) (msg string, err bool) {
	if m, _ := regexp.MatchString("TODO[ !:]", line); m {
		msg, err = line, true
	}
	return
}

// FixmeLint is a stateless linter that checks for and prints out lines that
// have "FIXM"E notices.
type FixmeLint struct{}

func (FixmeLint) Lint(line string) (msg string, err bool) {
	if m, _ := regexp.MatchString("FIXME[ !:]", line); m {
		msg, err = line, true
	}
	return
}

// XXXLint is a stateless linter that checks for and prints out lines that
// have XXX-style "TODO" notices (a convention in Java, at least).
type XXXLint struct{}

func (XXXLint) Lint(line string) (msg string, err bool) {
	if m, _ := regexp.MatchString("XXX[ !:]", line); m {
		msg, err = line, true
	}
	return
}
