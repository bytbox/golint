// Copyright 2011 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/token"
	"sync"
)

type VariableDeprecationLinter struct {
	LinterDesc
	packageName string
	varName     string
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
			lints <- ParsingLint{vdl, fset.Position(n.Pos())}
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
}

type MethodDeprecationLinter struct {
	LinterDesc
	packageName string
	typeName    string
	methodName  string
}

