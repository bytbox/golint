package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode"
)

const (
	T_COMMENT = iota
	T_SEMI
	T_IDENT
	T_KEYWORD
	T_OPER
	T_INT
	T_FLOAT
	T_IMAG
	T_CHAR
	T_STRING
)

type runeReader interface {
	ReadRune() (rune, int, error)
}

type Position struct {
	Line uint32
	Char uint32
}

type Token struct {
	Kind int
	Data interface{}
	Pos  Position
}

func (t Token) String() string {
	switch t.Kind {
	case T_IDENT: fallthrough
	case T_KEYWORD: fallthrough
	case T_OPER:
		return t.Data.(string)
	case T_SEMI:
		return "<semi>"
	case T_CHAR: fallthrough
	case T_STRING:
		return fmt.Sprintf("'%#v'", t.Data)
	}
	return "<tok>"
}

var keywords = map[string]bool{
	"break": true,
	"default": true,
	"func": true,
	"interface": true,
	"select": true,
	"case": true,
	"defer": true,
	"go": true,
	"map": true,
	"struct": true,
	"chan": true,
	"else": true,
	"goto": true,
	"package": true,
	"switch": true,
	"const": true,
	"fallthrough": true,
	"if": true,
	"range": true,
	"type": true,
	"continue": true,
	"for": true,
	"import": true,
	"return": true,
	"var": true,
}

func isKeyword(id string) bool {
	_, ok := keywords[id]
	return ok
}

// White space, formed from spaces (U+0020), horizontal tabs (U+0009), carriage
// returns (U+000D), and newlines (U+000A)
func isWS(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func LexString(src []byte) (ts []Token, err error) {
	var ltok *Token
	s := bytes.NewReader(src)
	r, _, err := s.ReadRune()

	pos := Position{1, 0}
	token := func(k int, d interface{}) {
		ts = append(ts, Token{k, d, pos})
		ltok = &ts[len(ts)-1]
	}
	nextR := func() bool {
		r, _, err = s.ReadRune()
		pos.Char++
		return err == nil
	}
	mustR := func() bool {
		if !nextR() {
			fmt.Errorf("unexpected error: %s\n", err.Error())
			return false
		}
		return true
	}

	for err == nil {
		if r == '\n' { goto Newline }
		if isWS(r) { goto Next }

		// identifiers and keywords
		if isLetter(r) {
			// identifier = letter { letter | unicode_digit } .
			str := string(r)
			for nextR() && (isLetter(r) || unicode.IsDigit(r)) {
				str += string(r)
			}
			if isKeyword(str) {
				token(T_KEYWORD, str)
				goto Next
			} else {
				token(T_IDENT, str)
				goto Next
			}
		}

		if r == ';' {
			token(T_SEMI, ";")
			goto Next
		}

		// character literal
		if r == '\'' {
			if !mustR() { break }
			if r == '\'' {
				err = fmt.Errorf("unexpected single quote\n")
				break
			}
			if r == '\\' {
				/*
				\a   U+0007 alert or bell
				\b   U+0008 backspace
				\f   U+000C form feed
				\n   U+000A line feed or newline
				\r   U+000D carriage return
				\t   U+0009 horizontal tab
				\v   U+000b vertical tab
				\\   U+005c backslash
				\'   U+0027 single quote
				\000
				\x..
				\u....
				\U........
				*/
				if !mustR() { break }
				switch r {
				case 'a': token(T_CHAR, '\a')
				case 'b': token(T_CHAR, '\b')
				case 'f': token(T_CHAR, '\f')
				case 'n': token(T_CHAR, '\n')
				case 'r': token(T_CHAR, '\r')
				case 't': token(T_CHAR, '\t')
				case 'v': token(T_CHAR, '\v')
				case '\\': token(T_CHAR, '\\')
				case '\'': token(T_CHAR, '\'')
				case 'x':
				case 'u':
				case 'U':
				case '0': fallthrough
				case '1': fallthrough
				case '2': fallthrough
				case '3': fallthrough
				case '4': fallthrough
				case '5': fallthrough
				case '6': fallthrough
				case '7':
				}
				goto Next
			}

			c := r
			if !mustR() { break }
			if r != '\'' {
				err = fmt.Errorf("expected single quote\n")
				break
			}
			token(T_CHAR, c)
			goto Next
		}

		// string literal

		// numerical literals

		// operators, delimiters, etc.
Newline:
		// attempt semicolon insertion
		if ltok != nil {
			if ltok.Kind == T_IDENT ||
				ltok.Kind == T_INT ||
				ltok.Kind == T_FLOAT ||
				ltok.Kind == T_IMAG ||
				ltok.Kind == T_CHAR ||
				ltok.Kind == T_STRING ||
				ltok.Kind == T_KEYWORD && (
					ltok.Data == "break" ||
					ltok.Data == "continue" ||
					ltok.Data == "fallthrough" ||
					ltok.Data == "return" ) ||
				ltok.Kind == T_OPER && (
					ltok.Data == "++" ||
					ltok.Data == "--" ||
					ltok.Data == "}" ||
					ltok.Data == ")" ||
					ltok.Data == "]") {
				token(T_SEMI, "")
			}
			ltok = nil
		}
		pos.Char = 0
		pos.Line++
Next:
		r, _, err = s.ReadRune()
		pos.Char++
	}
	if err == io.EOF { err = nil }
	return
}

func Lex(r io.Reader) (ts []Token, err error) {
	src, err := ioutil.ReadAll(r)
	if err != nil { return }
	return LexString(src)
}

func PrintLex(r io.Reader) {
	ts, err := Lex(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
	for _, t := range ts {
		fmt.Println(t)
	}
}
