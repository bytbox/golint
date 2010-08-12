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
