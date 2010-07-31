package main

import (
	"fmt"
	"go/ast"
)

// Error represents an error in the source code
type Error struct {
	msg string
	done bool
}

type validationVisitor struct {
	sync chan Error
}

func (v validationVisitor) Visit(node interface{}) ast.Visitor {
	// check the type
	if decl, ok := node.(ast.BadDecl); ok {
		msg := fmt.Sprintf("bad declaration at %d",
			decl.Position)
		v.sync <- Error{msg:msg}
	}
	if stmt, ok := node.(ast.BadStmt); ok {
		msg := fmt.Sprintf("bad statement at %d",
			stmt.Position)
		v.sync <- Error{msg:msg}
	}
	if expr, ok := node.(ast.BadExpr); ok {
		msg := fmt.Sprintf("bad expression at %d",
			expr.Position)
		v.sync <- Error{msg:msg}
	}
	return v
}

func (v validationVisitor) WalkOn(file *ast.File) {
	ast.Walk(v, file)
	v.sync <- Error{done: true}
}

// ValidParseLint is a parsing linter that checks to make sure the code
// parses properly, and emits more helpful messages than the compiler when
// it does not.
type ValidParseLint struct {
	sync chan Error
}
func (l *ValidParseLint) Init(file *ast.File) {
	// don't buffer, to force a rate limit
	l.sync = make(chan Error)
	visitor := &validationVisitor{l.sync}
	go visitor.WalkOn(file)
}
func (l *ValidParseLint) Next() (msg string, err bool) {
	error := <- l.sync
	if error.done {
		return
	}
	msg = error.msg
	err = true
	return
}
