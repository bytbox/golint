// Copyright 2010 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
)

// FuncDeprecationLint is a parsing linter that looks for and reports usages
// of deprecated functions.
type FuncDeprecationLint struct {
	Name string // the name of the deprecated function, including the package name, if appropriate
	Reason string // the reason for deprecation, or an alternate form

	errors chan deprecatedFuncUse
}

type funcDeprecationVisitor struct {
	errors chan deprecatedFuncUse
}

type deprecatedFuncUse struct {
	lineno int // a lineno of -1 indicates that we're done
}

func FuncDep(name string, reason string) (lint *FuncDeprecationLint) {
	lint = &FuncDeprecationLint{}
	lint.Name, lint.Reason = name, reason
	return
}

func (v *funcDeprecationVisitor) WalkOn(file *ast.File) {
	ast.Walk(v, file)
	v.errors <- deprecatedFuncUse{lineno: -1}
}

func (v *funcDeprecationVisitor) Visit(node interface{}) ast.Visitor {	
	if ce, ok := node.(*ast.CallExpr); ok {
		fun := ce.Fun
		if id, ok := fun.(*ast.Ident); ok {
			// it's a function!
			fmt.Printf("%v\n", id);
		}
	}
	return v;
}

func (l *FuncDeprecationLint) Init(file *ast.File) {
	l.errors = make(chan deprecatedFuncUse)
	v := &funcDeprecationVisitor{}
	v.errors = l.errors
	// start the parse
	go v.WalkOn(file)
}

func (l *FuncDeprecationLint) Next() (msg string, err bool) {
	error := <-l.errors
	err = !(error.lineno==-1)
	msg = fmt.Sprintf("message")
	return
}

// MethodDeprecationLint is a parsing linter that looks for and reports usages
// of deprecated methods.
type MethodDeprecationLint struct{}
