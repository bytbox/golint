// Copyright 2010 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"strings"
)

// LineLengthLint is a stateless lint that checks that line lengths are
// reasonable.
type LineLengthLint struct{}

func (lint LineLengthLint) Lint(line string) (msg string, err bool) {
	length := 0
	// count characters - a tab is eight characters
	chars := strings.Split(line, "", -1)
	for _, c := range chars {
		if c == "\t" {
			length += 8
		} else {
			length++
		}
	}
	// limit is 78
	limit := 78
	if length > limit {
		msg = fmt.Sprintf("line too long (%d > %d)", length, limit)
		err = true
	}
	return
}

// TabsOnlyLint is a stateless lint that checks that only tabs are used to
// indent lines.
type TabsOnlyLint struct{}

func (lint TabsOnlyLint) Lint(line string) (msg string, err bool) {
	chars := strings.Split(line, "", -1)
	for _, c := range chars {
		if c == " " {
			msg = "spaces used for indentation"
			err = true
			break
		}
		if c != "\t" {
			break
		}
	}
	return
}

// TrailingWhitespaceLint is a stateless lint that checks that there is no
// trailing whitespace.
type TrailingWhitespaceLint struct{}

func (lint TrailingWhitespaceLint) Lint(line string) (msg string, err bool) {
	if len(line) == 0 {
		// it's a blank line, just return no error
		return
	}
	chars := strings.Split(line, "", -1)
	c := chars[len(chars)-1]
	if c == " " || c == "\t" {
		msg = "trailing whitespace"
		err = true
	}
	return
}

// SemicolonLint is a stateless lint that checks that there are no unneeded
// semicolons.
type SemicolonLint struct{}

func (lint SemicolonLint) Lint(line string) (msg string, err bool) {
	if len(line) == 0 {
		// it's a blank line, just return no error
		return
	}
	chars := strings.Split(line, "", -1)
	c := chars[len(chars)-1]
	if c == ";" {
		msg = "unnecesary semicolon"
		err = true
	}
	return
}

// FilesizeLint is a stateful lint that checks that the number of lines in
// a file is reasonable.
type FilesizeLint struct {
	linecount int
}

const lineLimit = 1200

func (l *FilesizeLint) Reset() {
	l.linecount = 0
}
func (l *FilesizeLint) Lint(line string, lineno int) (msg string, err bool) {
	l.linecount++
	return
}
func (l *FilesizeLint) Done() (msg string, err bool) {
	if l.linecount > lineLimit {
		msg = fmt.Sprintf("file too long: %d lines (%d max)",
			l.linecount, lineLimit)
		err = true
	}
	return
}

// TrailingNewlineLint is a stateful lint that checks that there is only one
// blank line at the end of the file.
type TrailingNewlineLint struct {
	blankLineCount int
}

func (l *TrailingNewlineLint) Reset() {
	l.blankLineCount = 0
}
func (l *TrailingNewlineLint) Lint(line string, lineno int) (msg string, err bool) {
	if len(line) == 0 {
		l.blankLineCount++
	} else {
		l.blankLineCount = 0
	}
	return
}
func (l *TrailingNewlineLint) Done() (msg string, err bool) {
	if l.blankLineCount > 1 {
		msg = "extra trailing blank lines (only one permitted)"
		err = true
	}
	if l.blankLineCount < 1 {
		msg = "no trailing blank line"
		err = true
	}
	return
}

// UncleanImportLint is a parsing linter that looks for and reports ugly
// package imports.
type UncleanImportLint struct {
	complaints chan uncleanImport
}

type uncleanImport struct {
	path string // the given path too the package
	name string // the name the package was imported as
	done bool   // are we done?
}

func (l *UncleanImportLint) Init(file *ast.File) {
	l.complaints = make(chan uncleanImport)
	v := &uncleanImportVisitor{}
	v.complaints = l.complaints
	// start the parse
	go v.WalkOn(file)
}

type uncleanImportVisitor struct {
	complaints chan uncleanImport
}

func (v *uncleanImportVisitor) WalkOn(file *ast.File) {
	ast.Walk(v,file)
	v.complaints <- uncleanImport{done: true}
}

func (v *uncleanImportVisitor) Visit(node interface{}) ast.Visitor {
	if is, ok := node.(*ast.ImportSpec); ok {
		path := strings.Trim(string(is.Path.Value), "\"")
		name := is.Name
		if name != nil { // if name is nil, it's clean
			namestr := name.Name
			if namestr == "." { // XXX check for other forms of clobbering
				// this will clobber stuff - complain
				v.complaints <- uncleanImport{path, namestr, false}
			}
		}
	}
	return v
}

func (l *UncleanImportLint) Next() (msg string, err bool) {
	complaint := <- l.complaints
	err = !complaint.done
	msg = fmt.Sprintf("unclean import of %s as %s (may clobber)",
			complaint.path, complaint.name)
	return
}
