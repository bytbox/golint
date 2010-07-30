.PHONY: all clean

include ${GOROOT}/src/Make.${GOARCH}

all: golint

golint: main.${O}
	${LD} -o $@ main.${O}

MAINFILES = main.go style.go valid.go

main.6: ${MAINFILES}
	${GC} -o $@ ${MAINFILES}

clean:
	rm -f golint *.${O}
