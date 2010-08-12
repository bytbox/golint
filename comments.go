// Copyright 2010 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"regexp"
	"strings"
)

// TodoLint is a stateless linter that checks for and prints out lines that
// have "TODO" notices.
type TodoLint struct{}

func (TodoLint) Lint(line string) (msg string, err bool) {
	if m, _ := regexp.MatchString("TODO[ !:]", line); m {
		msg, err = strings.Trim(line, "\t"), true
	}
	return
}

// FixmeLint is a stateless linter that checks for and prints out lines that
// have "FIXM"E notices.
type FixmeLint struct{}

func (FixmeLint) Lint(line string) (msg string, err bool) {
	if m, _ := regexp.MatchString("FIXME[ !:]", line); m {
		msg, err = strings.Trim(line, "\t"), true
	}
	return
}

// XXXLint is a stateless linter that checks for and prints out lines that
// have XXX-style "TODO" notices (a convention in Java, at least).
type XXXLint struct{}

func (XXXLint) Lint(line string) (msg string, err bool) {
	if m, _ := regexp.MatchString("XXX[ !:]", line); m {
		msg, err = strings.Trim(line, "\t"), true
	}
	return
}
