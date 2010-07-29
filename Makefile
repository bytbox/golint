.PHONY: all clean

all: golint

golint: main.6
	6l -o $@ main.6

main.6: main.go style.go
	6g -o $@ $?

clean:
	rm golint *.6
