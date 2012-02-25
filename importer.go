package main

import (
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
	pkgname := path.Base(pathstr)
	rp, ok := imports[pkgname]
	if ok { return rp, nil }
	pr := filepath.Join(runtime.GOROOT(), "src", "pkg")
	pp := filepath.Join(pr, pathstr)
	pkg, err := ast.NewPackage(token.NewFileSet(), pkglist, importer, universe)
	pkgs, err := parser.ParseDir(token.NewFileSet(), pp, isGoSource, 0)
	if err != nil {
		return nil, err
	}
	pkgObj := ast.NewObj(ast.Pkg, pkgname)
	pkgObj.Data = pkg.Scope
	imports[pkgname] = pkgObj
	return pkgObj, nil
}

