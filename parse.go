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

// Conversion between ast.Visitor and the simpler type func(ast.Node) (bool).
type SimpleVisitor struct {
	visitFunc func(node ast.Node) bool
}

// Returns itself if 'visitFunc' returns true - otherwise returns nil.
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
			errs chan os.Error) {
	lintWG := new(sync.WaitGroup)
	lintWG.Add(1)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		errs <- err
	}

	// validate the parse
	valid := true
	ast.Walk(visitor(func (node ast.Node) bool {
		switch t := node.(type) {
		case *ast.BadDecl:
		case *ast.BadExpr:
		case *ast.BadStmt:
			valid = false
		}
		return true
	}), file)
	if !valid {
		println("not valid!")
		return
	}

	for _, linter := range ParsingLinters {
		lintWG.Add(1)
		go func() {
			nodeChan := make(chan ast.Node)
			go linter.RunLint(nodeChan, lintRoot, lintWG)

			ast.Walk(visitor(func (node ast.Node) bool {
				return true
			}), file)
			lintWG.Done()
		}()
	}

	lintWG.Done()
	lintWG.Wait()
}

// Represents a parse-based linter.
//
// The linter is given a stream of nodes in an AST (which is guaranteed to be
// valid, with no BadX nodes).
type ParsingLinter interface {
	String() string
	RunLint(chan ast.Node, chan Lint, *sync.WaitGroup)
}

