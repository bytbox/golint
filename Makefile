# Copyright 2010 The Golint Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY: all clean install

include ${GOROOT}/src/Make.inc

all: golint

golint: main.${O}
	${LD} -o $@ main.${O}

MAINFILES = main.go \
            line.go \
            rules.go \
            util.go
            

main.${O}: ${MAINFILES}
	${GC} -o $@ ${MAINFILES}

rules.go: genrules.pl $(wildcard rules/*/* rules/*)
	perl genrules.pl

install: /usr/local/bin/golint

/usr/local/bin/golint: golint
	cp $? $@

clean:
	rm -f golint *.${O} rules.go

