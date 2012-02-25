package main

import (
	. "go/ast"
)

var (
	universe = &Scope{
		Outer:   nil,
		Objects: map[string]*Object{
			"string": NewObj(Typ, "string"),
		},
	}
)

