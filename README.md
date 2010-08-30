Golint is a linter for golang, providing basic style and syntax checking as
well as scanning for common errors.


Gofmt Compatibility
-------------------

Golint is not compatible with the current version of gofmt (the main go code
formatting program), in the sense that running gofmt on code that passes
golint may produce code that does not pass golint. This is the product of a
combination of bugs in both gofmt and the go compiler, and should be fixed
soon.

In the long term, golint hopes to be compatible with gofmt in the sense above,
although running gofmt on code that does not completely pass golint may
increase the number of errors, and golint will never object to all code that
gofmt would change.


Dependencies
------------

Golint uses opts for option parsing, available here:
<http://opts-go.googlecode.com/>. You can install it automatically with
`goinstall`:

    goinstall -u opts-go.googlecode.com/hg


Bugs
----

Please report bugs and feature requests to the issue tracker on github:
<http://github.com/bytbox/golint/issues>.
