// Copyright 2011 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"container/vector"
	"io/ioutil"
	"strings"
	"os"
	"path"
)

// Misc. utility functions

// Read an array of lines from a file
func ReadFileLines(filename string) ([]string, os.Error) {
	contents, err := ioutil.ReadFile(filename)
	lines := strings.Split(string(contents), "\n")
	return lines, err
}

// Expand a list of files and directories
func ExpandFiles(files []string) []string {
	flatFiles := new(vector.StringVector)
	for _, filename := range files {
		info, err := os.Stat(filename)
		if err != nil {
			continue
		}
		if info.IsDirectory() {
			flatFiles.AppendVector(listFiles(filename))
		} else {
			flatFiles.Push(filename)
		}
	}
	return *flatFiles
}

// Extract all strings with the given suffix.
func FilterSuffix(suffix string, strs []string) []string {
	newstrs := make([]string, len(strs))
	i := 0
	for _, f := range strs {
		if strings.HasSuffix(f, suffix) {
			newstrs[i] = f
			i = i + 1
		}
	}
	return newstrs[0:i]
}

// Call two nullary functions in order.
func Seq(a func(), b func()) {
	a()
	b()
}

// List the names of all regular files in a given directory.
func listFiles(fname string) *vector.StringVector {
	files := new(vector.StringVector)
	f, _ := os.Open(fname)
	dn, _ := f.Readdirnames(-1)
	for _, filename := range dn {
		if filename[0] == '.' {
			continue
		}
		filename = path.Join(fname, filename)
		info, err := os.Stat(filename)
		if err != nil {
			continue
		}
		if info.IsDirectory() {
			files.AppendVector(listFiles(filename))
		} else if info.IsRegular() {
			files.Push(filename)
		}
	}
	return files
}

