// Copyright 2011 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"sync"
)

type DeprecationNotes struct {
	gofix string
}

func (dn DeprecationNotes) String() string {
	if len(dn.gofix)>0 {
		return fmt.Sprintf("fix with `gofix %s`", dn.gofix)
	}
	return ""
}

type VariableDeprecationLinter struct {
	LinterDesc
	packageName string
	varName     string
	extra       DeprecationNotes
}

func (vdl VariableDeprecationLinter) RunLint(
		fset *token.FileSet,
		nodes chan ast.Node,
		lints chan Lint,
		wg *sync.WaitGroup) {
	wg.Add(1)
	getVarRefs(nodes, func(n *ast.SelectorExpr, e *ast.Ident) {
		packageName := e.String() // TODO actually look up the package
		varName := n.Sel.String()
		if packageName == vdl.packageName &&
			varName == vdl.varName {
			lints <- ParsingLint{vdl,
				fset.Position(n.Pos()), vdl.extra.String()}
		}
	})
	wg.Done()
}

func getVarRefs(nodes chan ast.Node, f func(*ast.SelectorExpr, *ast.Ident)) {
	for node := range nodes {
		switch n := node.(type) {
		case (*ast.SelectorExpr):
			switch e := n.X.(type) {
			case (*ast.Ident):
				f(n, e)
			}
		}
	}
}

type FunctionDeprecationLinter struct {
	LinterDesc
	packageName string
	funcName    string
	args        []string
	extra       DeprecationNotes
}

func (fdl FunctionDeprecationLinter) RunLint(
		fset *token.FileSet,
		ns chan ast.Node,
		lints chan Lint,
		wg *sync.WaitGroup) {
	wg.Add(1)
	getFuncCalls(ns, func(pkg *ast.Ident, f *ast.Ident, args []ast.Expr) {
		// TODO actually look up the package
		if pkg.String() == fdl.packageName &&
			f.String() == fdl.funcName {
			lints <- ParsingLint{fdl,
				fset.Position(pkg.Pos()), fdl.extra.String()}
		}
	})
	wg.Done()
}

func getFuncCalls(nodes chan ast.Node,
		f func(*ast.Ident, *ast.Ident, []ast.Expr)) {
	for node := range nodes {
		n, ok := node.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pkg, ok := sel.X.(*ast.Ident)
		if !ok {
			continue
		}
		f(pkg, sel.Sel, n.Args)
	}
}

type MethodDeprecationLinter struct {
	LinterDesc
	packageName string
	typeName    string
	methodName  string
}

