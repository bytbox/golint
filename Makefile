# Copyright 2010 The Golint Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY: all clean install

include ${GOROOT}/src/Make.${GOARCH}

all: golint

golint: main.${O}
	${LD} -o $@ main.${O}

MAINFILES = main.go \
	data.go \
	style.go \
	valid.go \
	comments.go \
	deprecation.go \

main.${O}: ${MAINFILES}
	${GC} -o $@ ${MAINFILES}

install: /usr/local/bin/golint

/usr/local/bin/golint: golint
	cp $? $@

clean:
	rm -f golint *.${O}
