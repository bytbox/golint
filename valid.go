package main

import (
	"go/ast"
)

// ValidParseLint is a parsing linter that checks to make sure the code 
// parses properly, and emits more helpful messages than the compiler when
// it does not.
type ValidParseLint struct {
	file *ast.File
}
func (l ValidParseLint) Init(file *ast.File) {
	l.file = file
}
func (l ValidParseLint) Next() (msg string, err bool) {
	return
}
