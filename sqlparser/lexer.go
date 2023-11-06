package sqlparser

import (
	"io"
	"strings"
	"text/scanner"
)

type Lexer struct {
	s scanner.Scanner

	Value string
	Token rune
}

func NewLexer(rd io.Reader) *Lexer {
	l := &Lexer{}

	l.s = scanner.Scanner{}
	l.s.Mode = scanner.ScanIdents
	l.s.Init(rd)

	return l
}

func (l *Lexer) Lex() bool {
	l.Token = l.s.Scan()
	l.Value = l.s.TokenText()

	if l.Value == "<" {
		p := l.s.Peek()
		if p == '>' || p == '=' {
			l.Token = l.s.Scan()
			l.Value += l.s.TokenText()
		}
	} else if l.Value == ">" {
		p := l.s.Peek()
		if p == '=' {
			l.Token = l.s.Scan()
			l.Value += l.s.TokenText()
		}
	} else if strings.Contains(l.Value, "'") {
		l.Token = scanner.Ident
	} else if strings.Contains(l.Value, "\"") {
    l.Token = scanner.Ident
    l.Value = l.Value[1:len(l.Value)-1]
  }

	return l.Token != scanner.EOF
}
