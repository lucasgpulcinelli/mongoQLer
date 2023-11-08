package sqlparser

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

// The parser is implemented as a recursive descent parser. All rules are
// described in this file via comments in extended Backus-Naur form (mostly).

// Parse parses an SQL string and returns the statement that describes it.
//
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

// ParseBoolExpr parses a boolean expression only, returning the describing
// BooleanExpression.
//
// ParseBoolExpr -> BoolExpr
func ParseBoolExpr(sql string) (BooleanExpression, error) {
	l := NewLexer(strings.NewReader(sql))

	var be BooleanExpression

	if !l.Lex() {
		return nil, fmt.Errorf("failed parsing any SQL text")
	}

	if !BoolExpr(l, &be) {
		return nil, fmt.Errorf("failed parsing boolean expression")
	}

	if l.Lex() {
		return nil, fmt.Errorf("failed parsing SQL end: there is trailing input")
	}

	return be, nil

}

// SelectStmt -> <SELECT> Columns
// Columns -> ColumnOrGroup { <,> ColumnOrGroup } | <*>
func SelectStmt(l *Lexer, stmt *Statement) bool {
	if strings.ToUpper(l.Value) != "SELECT" || !l.Lex() {
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
		if strings.ToUpper(l.Value) == "FROM" {
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
	if strings.ToUpper(l.Value) != "FROM" || !l.Lex() {
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
	if strings.ToUpper(l.Value) != "JOIN" {
		return true
	}

	if !l.Lex() {
		return false
	}

	if l.Token != scanner.Ident {
		return false
	}

	stmt.JoinTable = l.Value

	if !l.Lex() || strings.ToUpper(l.Value) != "ON" || !l.Lex() {
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
	if strings.ToUpper(l.Value) != "WHERE" {
		stmt.Where = EmptyComparision{}
		return true
	}

	if !l.Lex() {
		return false
	}

	return BoolExpr(l, &stmt.Where)
}

// BoolExpr -> CompExpr { BoolOp CompExpr }
// BoolOp -> <AND> | <OR>
func BoolExpr(l *Lexer, be *BooleanExpression) bool {
	var comp BooleanExpression

	if !CompExpr(l, &comp) {
		return false
	}

	if l.Value == ";" || l.Value == ")" {
		*be = comp
		return true
	}

	bcomposite := &BooleanComposite{}
	bcomposite.SubExpr = make([]BooleanExpression, 1)
	bcomposite.SubExpr[0] = comp

	for {
		if bcomposite.BoolOp == "" {
			bcomposite.BoolOp = l.Value
		}

		if bcomposite.BoolOp != l.Value || !l.Lex() {
			return false
		}

		if !CompExpr(l, &comp) {
			return false
		}

		bcomposite.SubExpr = append(bcomposite.SubExpr, comp)

		if l.Value == ";" || l.Value == ")" {
			*be = bcomposite
			return true
		}
	}
}

// GetValue obtains an sql value from a string as an int, nil or string.
func GetValue(s string) any {
	if strings.ToUpper(s) == "NULL" {
		return nil
	}

	v, err := strconv.ParseInt(s, 0, 0)
	if err == nil {
		return v
	}

	if s[0] == '\'' || s[0] == '"' {
		return s[1 : len(s)-1]
	}

	return s
}

// SimpleCompExpr -> <ID> CompOp <ID> | <ID> <IS> (<NOT> | eps) <NULL>
func SimpleCompExpr(l *Lexer, be *BooleanExpression, id string) bool {
	comp := &Comparision{}
	comp.Id = id

	if strings.ToUpper(l.Value) == "IS" {
		if !l.Lex() {
			return false
		}

		if strings.ToUpper(l.Value) == "NOT" {
			comp.Op = "<>"
			if !l.Lex() {
				return false
			}
		} else {
			comp.Op = "="
		}

		if strings.ToUpper(l.Value) != "NULL" || !l.Lex() {
			return false
		}

		comp.Value = nil
	} else {
		comp.Op = l.Value
		if !l.Lex() || (l.Token != scanner.Ident && l.Token != scanner.Int) {
			return false
		}
		comp.Value = GetValue(l.Value)

		if !l.Lex() {
			return false
		}
	}

	*be = comp
	return true
}

// InCompExpr -> <ID> (<NOT> | eps) <IN> <(> ValueList <)>
// ValueList -> <ID> { <,> <ID> }
func InCompExpr(l *Lexer, be *BooleanExpression, id string) bool {
	incomp := &InComparision{}
	incomp.Id = id
	incomp.Values = make([]any, 1)

	if strings.ToUpper(l.Value) == "NOT" {
		incomp.Not = true
		if !l.Lex() {
			return false
		}
	}

	if strings.ToUpper(l.Value) != "IN" || !l.Lex() ||
		l.Value != "(" || !l.Lex() {
		return false
	}

	if !(l.Token == scanner.Ident || l.Token == scanner.Int) {
		return false
	}

	incomp.Values[0] = GetValue(l.Value)

	if !l.Lex() {
		return false
	}

	for l.Token == ',' {
		if !l.Lex() {
			return false
		}

		if !(l.Token == scanner.Ident || l.Token == scanner.Int) {
			return false
		}

		incomp.Values = append(incomp.Values, GetValue(l.Value))

		if !l.Lex() {
			return false
		}
	}

	if l.Value != ")" || !l.Lex() {
		return false
	}

	*be = incomp
	return true
}

// CompExpr -> <(> BoolExpr <)> | InCompExpr | SimpleCompExpr
func CompExpr(l *Lexer, be *BooleanExpression) bool {

	if l.Token == '(' {
		if !l.Lex() {
			return false
		}

		if !BoolExpr(l, be) {
			return false
		}

		if l.Token != ')' || !l.Lex() {
			return false
		}

		return true
	}

	if l.Token != scanner.Ident && l.Token != scanner.Int {
		return false
	}

	id := l.Value

	if !l.Lex() {
		return false
	}

	if l.Value == ">" || l.Value == ">=" || l.Value == "<" || l.Value == "<=" ||
		l.Value == "=" || l.Value == "<>" || l.Value == "IS" {
		return SimpleCompExpr(l, be, id)
	}

	if l.Value == "NOT" || l.Value == "IN" {
		return InCompExpr(l, be, id)
	}

	return false
}
