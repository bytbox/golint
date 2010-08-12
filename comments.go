package main

import "strings"

// TodoLint is a stateless linter that checks for and prints out lines that
// have TODO notices.
type TodoLint struct{}

func (TodoLint) Lint(line string) (msg string, err bool) {
	if strings.Index(line, "TODO") != -1 {
		msg, err = line, true
	}
	return
}

// FixmeLint is a stateless linter that checks for and prints out lines that
// have FIXME notices.
type FixmeLint struct{}

func (FixmeLint) Lint(line string) (msg string, err bool) {
	if strings.Index(line, "FIXME") != -1 {
		msg, err = line, true
	}
	return
}

// XXXLint is a stateless linter that checks for and prints out lines that
// have XXX-style TODO notices (a convention in Java, at least).
type XXXLint struct{}

func (XXXLint) Lint(line string) (msg string, err bool) {
	if strings.Index(line, "XXX") != -1 {
		msg, err = line, true
	}
	return
}
