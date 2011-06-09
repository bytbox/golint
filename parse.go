// Copyright 2011 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sync"
)

type SimpleVisitor struct {
	visitFunc func(node ast.Node) bool
}

func (v SimpleVisitor) Visit(node ast.Node) ast.Visitor {
	if v.visitFunc(node) {
		return v
	}
	return nil
}

func visitor(visit func(node ast.Node) bool) ast.Visitor {
	return SimpleVisitor{visit}
}

func RunParsingLinters(filename string,
			lintRoot chan Lint,
			errs chan os.Error,
			lintWG *sync.WaitGroup) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		errs <- err
	}

	ast.Walk(visitor(func (node ast.Node) bool {
		return false
	}), file)
}

// Represents a parse-based linter.
//
// The linter is given a stream of nodes in an AST (which is guaranteed to be
// valid, with no BadX nodes).
type ParsingLinter interface {
	String() string
	RunLint(chan ast.Node, chan Lint, *sync.WaitGroup)
}

