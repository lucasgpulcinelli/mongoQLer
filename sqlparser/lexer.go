package sqlparser

import (
	"io"
	"strings"
	"text/scanner"
)

// struct Lexer implements a simple lexer for the sqlparser using a text
// scanner.
type Lexer struct {
	s scanner.Scanner // the scanner itself

	Value string // the current value read as a string
	Token rune   // the current token mostly as returned from the scanner
}

// NewLexer creares a new Lexer from a reader.
func NewLexer(rd io.Reader) *Lexer {
	l := &Lexer{}

	l.s = scanner.Scanner{}
	l.s.Mode = scanner.ScanIdents
	l.s.Init(rd)

	l.s.Error = func(s *scanner.Scanner, msg string) {}

	return l
}

// Lex advances the scanner forward by one token, updating Token and Value
// data, and returning if the lexer reached EOF.
func (l *Lexer) Lex() bool {
	l.Token = l.s.Scan()
	l.Value = l.s.TokenText()

	// convert some tokens into more usable forms for SQL:
	if l.Value == "<" {
		// <= and <> are one token
		p := l.s.Peek()
		if p == '>' || p == '=' {
			l.Token = l.s.Scan()
			l.Value += l.s.TokenText()
		}
	} else if l.Value == ">" {
		// >= is one token
		p := l.s.Peek()
		if p == '=' {
			l.Token = l.s.Scan()
			l.Value += l.s.TokenText()
		}
	} else if strings.Contains(l.Value, "'") {
		// 'value' is an identifier token (for our uses), not a character
		l.Token = scanner.Ident
	} else if strings.Contains(l.Value, "\"") {
		// "value" is an identifier token, not a string, and the "" should be
		// removed because SQL ignores them
		l.Token = scanner.Ident
		l.Value = l.Value[1 : len(l.Value)-1]
	}

	return l.Token != scanner.EOF
}
