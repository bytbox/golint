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
	return DeprecationLint{fn+"()","[^a-zA-Z0-9_]"+fn+" *\\(", reason}
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
	err bool
}

func (v *packageDeprecationVisitor) Visit(node interface{}) ast.Visitor {
	if is, ok := node.(*ast.ImportSpec); ok {
		path := strings.Trim(string(is.Path.Value),"\"")
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
			l.Package, l.Reason), true
		l.err = false
	}
	return
}
