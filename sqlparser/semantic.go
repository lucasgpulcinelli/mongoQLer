package sqlparser

import (
	"fmt"
	"strings"

	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"go.mongodb.org/mongo-driver/bson"
)

// struct Column represents a parsed SQL column (select entry), which is either
// an identifier or an identifier with a group function associated. Note that
// naming columns is not supported.
type Column struct {
	Name          string
	GroupFunction string
}

// struct Statement represents a parsed SQL statement, with selection columns,
// a single origin table, a single (optional) joined table with a single join
// condition, and a filtering expression.
type Statement struct {
	SelectColumn []Column

	FromTable string

	JoinTable    string
	JoinFromAttr string
	JoinToAttr   string

	Where BooleanExpression
}

// A BooleanExpression represents a parsed boolean comparision that can be
// converted to a mongoDB bson document given the tables related in the query
// (for _id management).
type BooleanExpression interface {
	GetBson(tableFrom, tableJoin string) (bson.D, error)
}

// struct EmptyComparision represents a comparision that is always true
type EmptyComparision struct{}

// struct Comparision represents a simple comparision such as "A > 10", having
// an identifier on the left, a value on the right, and a boolean operator.
type Comparision struct {
	Id    string
	Value any
	Op    string
}

// struct InComparision represents a comparision using IN or NOT IN, such as
// "A IN (1, 2, 3, 4)".
type InComparision struct {
	Id     string
	Not    bool
	Values []any
}

// struct BooleanComposite represents a boolean expression with many sub
// boolean expressions and an operator joining them (such as AND or OR).
type BooleanComposite struct {
	BoolOp  string
	SubExpr []BooleanExpression
}

// IsAggregate returns if the Statement is an aggregation or a find.
// A Statement is an aggregation only if it has either a join or a group
// function in it.
func (stmt *Statement) IsAggregate() bool {
	if stmt.JoinTable != "" {
		return true
	}

	for _, col := range stmt.SelectColumn {
		if col.GroupFunction != "" {
			return true
		}
	}

	return false
}

// GetBson implements the BooleanExpression interface.
func (e EmptyComparision) GetBson(_, _ string) (bson.D, error) {
	return bson.D{}, nil
}

// GetBson implements the BooleanExpression interface.
func (c *Comparision) GetBson(tableFrom, tableJoin string) (bson.D, error) {
	operator := ""

	switch c.Op {
	default:
		return bson.D{}, fmt.Errorf("invalid operator %s", operator)
	case "=":
		operator = "$eq"
	case "<>":
		operator = "$ne"
	case ">":
		operator = "$gt"
	case ">=":
		operator = "$gte"
	case "<":
		operator = "$lt"
	case "<=":
		operator = "$lte"
	}

	var k string
	if oracleManager.TableContainsColumn(tableJoin, c.Id) {
		k = tableJoin + "." + keyManager.ToMongoId(tableJoin, c.Id)
	} else {
		k = keyManager.ToMongoId(tableFrom, c.Id)
	}

	return bson.D{{
		Key:   k,
		Value: bson.D{{Key: operator, Value: c.Value}},
	}}, nil

}

// GetBson implements the BooleanExpression interface.
func (ic *InComparision) GetBson(tableFrom, tableJoin string) (bson.D, error) {
	operator := "$in"
	if ic.Not {
		operator = "$nin"
	}

	var k string
	if oracleManager.TableContainsColumn(tableJoin, ic.Id) {
		k = tableJoin + "." + keyManager.ToMongoId(tableJoin, ic.Id)
	} else {
		k = keyManager.ToMongoId(tableFrom, ic.Id)
	}

	return bson.D{{
		Key:   k,
		Value: bson.D{{Key: operator, Value: ic.Values}},
	}}, nil
}

// GetBson implements the BooleanExpression interface.
func (bc *BooleanComposite) GetBson(
	tableFrom, tableJoin string,
) (bson.D, error) {

	boolOpStr := ""

	if strings.ToUpper(bc.BoolOp) == "AND" {
		boolOpStr = "$and"
	} else if strings.ToUpper(bc.BoolOp) == "OR" {
		boolOpStr = "$or"
	} else {
		return bson.D{}, fmt.Errorf("invalid boolean operator in WHERE")
	}

	sexprs := make([]bson.D, 0)
	for _, se := range bc.SubExpr {
		bs, err := se.GetBson(tableFrom, tableJoin)
		if err != nil {
			return bson.D{}, err
		}

		sexprs = append(sexprs, bs)
	}

	return bson.D{{Key: boolOpStr, Value: sexprs}}, nil
}
