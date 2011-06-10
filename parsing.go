// Copyright 2011 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
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
		return
	}

	for _, linter := range ParsingLinters {
		lintWG.Add(1)
		go func() {
			nodeChan := make(chan ast.Node)
			go linter.RunLint(filename, nodeChan, lintRoot, lintWG)

			ast.Walk(visitor(func (node ast.Node) bool {
				nodeChan <- node
				return true
			}), file)
			close(nodeChan)
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
	RunLint(string, chan ast.Node, chan Lint, *sync.WaitGroup)
}

type OverlappingImportsLinter struct {
}

func (oil OverlappingImportsLinter) Desc() LinterDesc {
	return LinterDesc{
		"misc",
		"overlapping-imports",
		"Imports of distinct packages should have distinct local names"}
}

func (oil OverlappingImportsLinter) String() string {
	return oil.Desc().String()
}

func (oil OverlappingImportsLinter) RunLint(
		filename string,
		nodes chan ast.Node,
		lints chan Lint,
		wg *sync.WaitGroup) {
	wg.Add(1)
	imports := make(map[string]string)
	for node := range nodes {
		switch node.(type) {
		case (*ast.ImportSpec):
			is := node.(*ast.ImportSpec)
			if is.Name != nil {
				imports[is.Path.Value] = is.Name.String()
			} else {
				path := strings.Trim(is.Path.Value, "\"")
				parts := strings.Split(path, "/", -1)
				imports[is.Path.Value] = parts[len(parts)-1]
			}
		}
	}

	localNameCount := make(map[string]int)
	for _, localName := range imports {
		localNameCount[localName] += 1
	}

	for localName, count := range localNameCount {
		if localName == "." {
			count += 1
		}
		if count > 1 {
			lints <- OverlappingImportsLint{oil,
				filename, localName, count}
		}
	}
	wg.Done()
}

type OverlappingImportsLint struct {
	linter Linter
	filename string
	localName string
	count int
}

func (oil OverlappingImportsLint) String() string {
	return fmt.Sprintf("%s in %s: '%s' used %d times",
		oil.linter.String(),
		oil.filename,
		oil.localName,
		oil.count)
}

