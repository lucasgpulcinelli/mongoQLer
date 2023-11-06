package sqlparser

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type Column struct {
	Name          string
	GroupFunction string
}

type Statement struct {
	SelectColumn []Column

	FromTable string

	JoinTable    string
	JoinFromAttr string
	JoinToAttr   string

	Where BooleanExpression
}

type BooleanExpression interface {
	GetBson() (bson.D, error)
}

type EmptyComparision struct{}

type Comparision struct {
	Id    string
	Value any
	Op    string
}

type InComparision struct {
	Id     string
	Not    bool
	Values []any
}

type BooleanComposite struct {
	BoolOp  string
	SubExpr []BooleanExpression
}

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

func (e EmptyComparision) GetBson() (bson.D, error) {
	return bson.D{}, nil
}

func (c *Comparision) GetBson() (bson.D, error) {
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

	return bson.D{{c.Id, bson.D{{operator, c.Value}}}}, nil

}

func (ic *InComparision) GetBson() (bson.D, error) {
	operator := "$in"
	if ic.Not {
		operator = "$nin"
	}

	return bson.D{{ic.Id, bson.D{{operator, ic.Values}}}}, nil
}

func (bc *BooleanComposite) GetBson() (bson.D, error) {
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
		bs, err := se.GetBson()
		if err != nil {
			return bson.D{}, err
		}

		sexprs = append(sexprs, bs)
	}

	return bson.D{{boolOpStr, sexprs}}, nil
}
