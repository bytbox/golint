# Copyright 2010 The Golint Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include ${GOROOT}/src/Make.inc

TARG = golint
GOLINTFILES = line.go parsing.go deprecation.go rules.go
GOFILES = main.go util.go $(addprefix go/lint/,${GOLINTFILES})

include $(GOROOT)/src/Make.cmd

go/lint/rules.go: genrules.pl $(wildcard rules/*/* rules/*)
	perl genrules.pl

