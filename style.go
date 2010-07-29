package main

import (
	"fmt"
	"strings"
)

// LineLengthLint is a stateless lint that checks that line lengths are 
// reasonable.
type LineLengthLint struct {}
func (lint LineLengthLint) Lint(line string) (msg string, err bool) {
	if line == "" {
		return
	}
 	length := 0
	// count characters - a tab is eight characters
	chars := strings.Split(line, "", -1)
	for _, c := range chars {
		if c == "\t" {
			length += 8
		} else {
			length++
		}
	}
	// limit is 78
	limit := 78
	if length > limit {
		msg = fmt.Sprintf("line too long (%d > %d)", length, limit)
		err = true
	}
	return
}

type FilesizeLint struct {}
