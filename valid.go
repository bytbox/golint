package main

import (
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
	err := Error{}
	// send whatever error along
	v.sync <- err
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
