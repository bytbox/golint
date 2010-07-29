.PHONY: all clean

all: golint

golint: main.6
	6l -o $@ main.6

MAINFILES = main.go style.go

main.6: ${MAINFILES}
	6g -o $@ ${MAINFILES}

clean:
	rm -f golint *.6
