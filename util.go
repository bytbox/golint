package main

import (
	"os"
	"strings"
)

func isGoSource(fi os.FileInfo) bool {
	return strings.HasSuffix(fi.Name(), ".go")
}
