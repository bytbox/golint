package main

import (
	"io/ioutil"
	"strings"
	"os"
)

// Misc. utility functions

// Read an array of lines from a file
func ReadFileLines(filename string) ([]string, os.Error) {
	contents, err := ioutil.ReadFile(filename)
	lines := strings.Split(string(contents), "\n", -1)
	return lines, err
}

// Expand a list of files and directories
func ExpandFiles(files []string) []string {
	return files // TODO
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

