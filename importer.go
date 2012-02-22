package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"
	"runtime"
)

var importer = ast.Importer(find_import)

// Implementing the ast.Importer signature
func find_import(imports map[string]*ast.Object, pathstr string) (*ast.Object, error) {
	pr := filepath.Join(runtime.GOROOT(), "src", "pkg")
	pp := filepath.Join(pr, pathstr)
	pkgs, err := parser.ParseDir(token.NewFileSet(), pp, isGoSource, 0)
	if err != nil {
		return nil, err
	}
	pkgname := path.Base(pathstr)
	pkg, ok := pkgs[pkgname]
	if !ok { return nil, errors.New("unknown") }
	pkgObj := ast.NewObj(ast.Pkg, pkgname)
	pkgObj.Data = pkg
	return pkgObj, nil
}

