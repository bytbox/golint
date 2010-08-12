// Copyright 2010 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
)

// DeprecationLint is a stateless linter that looks for and reports deprecated
// constructs.
//
// Not all deprecated forms may be expressable with a DeprecationLint.
type DeprecationLint struct {
	Form   string // a human-readable string representing the form
	Regexp string // a regular expression representing the form
	Reason string // the reason for deprecation, or an alternative form
}

func FuncDeprecationLint(fn string, reason string) (lint DeprecationLint) {
	return DeprecationLint{fn + "()",
		"[^a-zA-Z0-9_]" + regexp.QuoteMeta(fn) + " *\\(", reason}
}


func (l DeprecationLint) Lint(line string) (msg string, err bool) {
	if m, _ := regexp.MatchString(l.Regexp, line); m {
		err = true
		// XXX once the regexp patch is out, we can give more info
		msg = fmt.Sprintf("deprecated use of %s: %s",
			l.Form, l.Reason)
	}
	return
}

// MethodDeprecationLint is a parsing linter that looks for and reports usage
// of deprecated methods.
type MethodDeprecationLint struct {
	Type   string // the name of the type on which the method is called
	Method string // the name of the method
	Reason string // the reason for deprecation, or an alternative form
	file   *ast.File
}

func (l *MethodDeprecationLint) Init(file *ast.File) {
	l.file = file
	// TODO start methodDeprecationVisitor-based parse
}

type methodDeprecationVisitor struct{}

func (v *methodDeprecationVisitor) Visit(node interface{}) ast.Visitor {
	return v
}

func (l *MethodDeprecationLint) Next() (msg string, err bool) {
	// TODO use results of parse
	return
}

// PackageDeprecationLint is a parsing linter that looks for and reports usage
// of deprecated packages.
type PackageDeprecationLint struct {
	Package string // the name of the deprecated package
	Reason  string // the reason for deprecation, or an alternative form
	err     bool
}

func (l *PackageDeprecationLint) Init(file *ast.File) {
	visitor := &packageDeprecationVisitor{l.Package, false}
	ast.Walk(visitor, file)
	if visitor.err {
		l.err = true
	}
}

type packageDeprecationVisitor struct {
	pname string
	err   bool
}

func (v *packageDeprecationVisitor) Visit(node interface{}) ast.Visitor {
	if is, ok := node.(*ast.ImportSpec); ok {
		path := strings.Trim(string(is.Path.Value), "\"")
		if path == v.pname {
			v.err = true
			return nil
		}
	}
	return v
}

func (l *PackageDeprecationLint) Next() (msg string, err bool) {
	if l.err {
		msg, err = fmt.Sprintf(
			"use of deprecated package %s (%s)",
			l.Package, l.Reason),
			true
		l.err = false
	}
	return
}
