package main

// TODO catch places a semicolon would be inserted to cause a syntax error
// TODO unused warnings

import (
	"fmt"
	"go/ast"
	"go/parser"
	"io"
	"io/ioutil"
	"opts"
	"os"
	"strings"
)

var version = "0.0.4"

// options
var (
	disabledLintList = opts.Shortmulti("d", "disabled lints", "lint")
	disabledCatList  = opts.Shortmulti("D", "disabled lint categories", "category")
	showVersion      = opts.Longflag("version", "display version information and exit")
	verbose          = opts.Flag("v", "verbose", "print verbose output to standard error")
)

func main() {
	// Do the argument parsing
	opts.Usage("[sourcefile...]")
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
	// disable, via deletion, the disabled lints
	for _, dcat := range *disabledCatList {
		for lname := range statelessLinters {
			if category(lname)==dcat {
				disabledLintList.Push(lname)
			}
		}
		for lname := range statefulLinters {
			if category(lname)==dcat {
				disabledLintList.Push(lname)
			}
		}
		for lname := range parsingLinters {
			if category(lname)==dcat {
				disabledLintList.Push(lname)
			}
		}
	}
	for _, dlint := range *disabledLintList {
		statelessLinters[dlint] = nil, false
		statefulLinters[dlint] = nil, false
		parsingLinters[dlint] = nil, false
	}
	if *verbose {
		fmt.Fprintf(os.Stderr,
			"Beginning lint with %d stateless, "+
				"%d stateful, and %d parsing linters\n",
			len(statelessLinters),
			len(statefulLinters),
			len(parsingLinters))
	}
	for _, filename := range opts.Args {
		err := DoLintFrom(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"golint: couldn't lint file: %s\n",
				filename)
		}
	}
	if len(opts.Args) == 0 {
		// read from standard input
		DoLint(os.Stdin, "stdin")
	}
}

// Show version information
func ShowVersion() {
	fmt.Printf("golint v%s\n", version)
}

// category returns the category of a lint name
func category(lname string) string {
	parts := strings.Split(lname, ":", 2)
	return parts[0]
}

var statelessLinters = map[string]StatelessLinter{
	"style:linelen":            &LineLengthLint{},
	"style:tabsonly":           &TabsOnlyLint{},
	"style:trailingwhitespace": &TrailingWhitespaceLint{},
	"style:semicolon":          &SemicolonLint{},
	"deprecated:new":           FuncDeprecationLint("new", "use &T{}"),
	"todo:todo":                &TodoLint{},
	"todo:fixme":               &FixmeLint{},
	"todo:xxx":                 &XXXLint{},
}

var statefulLinters = map[string]StatefulLinter{
	"style:filesize":        &FilesizeLint{},
	"style:trailingnewline": &TrailingNewlineLint{},
}

var parsingLinters = map[string]ParsingLinter{
	"syntaxvalidparse": &ValidParseLint{},
	"deprecated:once": &PackageDeprecationLint{Package: "once",
		Reason: "use sync.Once"},
}

type StatelessLinter interface {
	Lint(string) (string, bool)
}

type StatefulLinter interface {
	Lint(string, int) (string, bool)
	Reset()
	Done() (string, bool)
}

type ParsingLinter interface {
	Init(*ast.File)
	Next() (string, bool)
}

func DoLintFrom(filename string) os.Error {
	// read in the file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// create a reader on the content, and call DoLint
	DoLint(strings.NewReader(string(content)), filename)
	return nil
}

// a parseResult represents the result of a call to ParseFile
type parseResult struct {
	file *ast.File
	err  os.Error
}

func parse(filename string, content string, c chan parseResult) {
	file, err := parser.ParseFile(filename, content, nil, 0)
	res := parseResult{file, err}
	c <- res
}

func DoLint(reader io.Reader, filename string) os.Error {
	// read in the file
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	// prepare all the stateful linters
	for _, linter := range statefulLinters {
		linter.Reset()
	}
	// start parsing in parallel
	c := make(chan parseResult)
	go parse(filename, string(content), c)
	// for each line
	lines := strings.Split(string(content), "\n", -1)
	for lineno, line := range lines {
		// run through the stateless linters
		for lname, linter := range statelessLinters {
			msg, err := linter.Lint(line)
			if err {
				fmt.Printf("%s:%d: %s (%s)\n",
					filename, lineno+1, msg, lname)
			}
		}
		// run through the stateful linters
		for lname, linter := range statefulLinters {
			msg, err := linter.Lint(line, lineno)
			if err {
				fmt.Printf("%s: %s (%s)\n",
					filename, msg, lname)
			}
		}
	}
	// tell all the stateful linters we're done
	for lname, linter := range statefulLinters {
		msg, err := linter.Done()
		if err {
			fmt.Printf("%s: %s (%s)\n", filename, msg, lname)
		}
	}
	// run the parsing linters
	// First, get the result of the parsing
	result := <-c
	if result.err != nil {
		fmt.Printf("%s (in parser)\n", result.err)
	}
	astFile := result.file
	// for each parsingLinter
	for lname, linter := range parsingLinters {
		linter.Init(astFile)
		msg, cont := linter.Next()
		for cont {
			fmt.Printf("%s: %s (%s)\n", filename, msg, lname)
			msg, cont = linter.Next()
		}
	}
	return nil
}
