package sqlparser

import (
	"fmt"
	"strings"
	"text/scanner"
)

// Parse -> SelectStmt FromStmt OptJoinStmt OptWhereStmt
func Parse(sql string) (*Statement, error) {
	l := NewLexer(strings.NewReader(sql))

	stmt := &Statement{}

	if !l.Lex() {
		return nil, fmt.Errorf("failed parsing any SQL text")
	}

	if !SelectStmt(l, stmt) {
		return nil, fmt.Errorf("failed parsing SQL SELECT")
	}

	if !FromStmt(l, stmt) {
		return nil, fmt.Errorf("failed parsing SQL FROM")
	}

	if !OptJoinStmt(l, stmt) {
		return nil, fmt.Errorf("failed parsing SQL JOIN")
	}

	if !OptWhereStmt(l, stmt) {
		return nil, fmt.Errorf("failed parsing SQL WHERE")
	}

	if l.Lex() {
		return nil, fmt.Errorf("failed parsing SQL end: there is trailing input")
	}

	return stmt, nil
}

// SelectStmt -> <SELECT> Columns
// Columns -> ColumnOrGroup { <,> ColumnOrGroup } | <*>
func SelectStmt(l *Lexer, stmt *Statement) bool {
	if l.Value != "SELECT" || !l.Lex() {
		return false
	}

	stmt.SelectColumn = make([]Column, 0)
	if l.Value == "*" {
		return l.Lex()
	}

	if !ColumnOrGroup(l, stmt) {
		return false
	}

	for {
		if l.Value == "FROM" {
			return true
		}

		if l.Value != "," || !l.Lex() {
			return false
		}

		if !ColumnOrGroup(l, stmt) {
			return false
		}
	}
}

// ColumnOrGroup -> <ID> OptGroup
// OptGroup -> <(> <ID> <)> | eps
func ColumnOrGroup(l *Lexer, stmt *Statement) bool {
	s := l.Value
	if l.Token != scanner.Ident || !l.Lex() {
		return false
	}

	if l.Value != "(" {
		stmt.SelectColumn = append(stmt.SelectColumn, Column{Name: s})
		return true
	}

	if !l.Lex() {
		return false
	}

	s2 := l.Value

	if !l.Lex() || l.Value != ")" || !l.Lex() {
		return false
	}

	stmt.SelectColumn = append(
		stmt.SelectColumn,
		Column{Name: s2, GroupFunction: s},
	)
	return true
}

// FromStmt -> <FROM> <ID>
func FromStmt(l *Lexer, stmt *Statement) bool {
	if l.Value != "FROM" || !l.Lex() {
		return false
	}

	if l.Token != scanner.Ident {
		return false
	}

	stmt.FromTable = l.Value

	return l.Lex()
}

// OptJoinStmt -> <JOIN> <ID> <ON> <ID> <=> <ID> | eps
func OptJoinStmt(l *Lexer, stmt *Statement) bool {
	if l.Value != "JOIN" {
		return true
	}

	if !l.Lex() {
		return false
	}

	if l.Token != scanner.Ident {
		return false
	}

	stmt.JoinTable = l.Value

	if !l.Lex() || l.Value != "ON" || !l.Lex() {
		return false
	}

	stmt.JoinFromAttr = l.Value

	if !l.Lex() || l.Value != "=" || !l.Lex() {
		return false
	}

	stmt.JoinToAttr = l.Value

	return l.Lex()
}

// OptWhereStmt -> <WHERE> BoolExpr | eps
func OptWhereStmt(l *Lexer, stmt *Statement) bool {
	if l.Value != "WHERE" {
		return true
	}

	if !l.Lex() {
		return false
	}

	stmt.Where = make([]Comparision, 0)
	return BoolExpr(l, stmt)
}

// BoolExpr -> CompExpr { BoolOp CompExpr }
// BoolOp -> <AND> | <OR>
func BoolExpr(l *Lexer, stmt *Statement) bool {
	if !CompExpr(l, stmt) {
		return false
	}

	for {
		if l.Value == ";" {
			return true
		}

		if stmt.BooleanOp == "" {
			stmt.BooleanOp = l.Value
		}

		if stmt.BooleanOp != l.Value || !l.Lex() {
			return false
		}

		if !CompExpr(l, stmt) {
			return false
		}
	}
}

// CompExpr -> <ID> CompOp <ID>
func CompExpr(l *Lexer, stmt *Statement) bool {
	c := Comparision{}

	if l.Token != scanner.Ident && l.Token != scanner.Int {
		return false
	}

	c.Left = l.Value

	if !l.Lex() {
		return false
	}

	if !(l.Value == ">" || l.Value == ">=" || l.Value == "<" ||
		l.Value == "<=" || l.Value == "=" || l.Value == "<>") {
		return false
	}

	c.Op = l.Value

	if !l.Lex() {
		return false
	}

	if l.Token != scanner.Ident && l.Token != scanner.Int {
		return false
	}

	c.Right = l.Value

	if !l.Lex() {
		return false
	}

	stmt.Where = append(stmt.Where, c)

	return true
}
