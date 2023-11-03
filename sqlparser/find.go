package sqlparser

import (
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func (stmt *Statement) GetFindSelect() (bson.D, error) {
	if len(stmt.SelectColumn) == 0 {
		return bson.D{}, nil
	}

	ret := bson.D{}

	for _, selection := range stmt.SelectColumn {
		ret = append(ret, bson.E{selection.Name, 1})
	}

	return ret, nil
}

func (stmt *Statement) GetFindWhere(i int) (string, bson.M, error) {
	comp := stmt.Where[i]

	operator := ""

	switch comp.Op {
	default:
		return "", bson.M{}, fmt.Errorf("invalid operator %s", operator)
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

	v, err := strconv.ParseInt(comp.Right, 0, 0)
	if err == nil {
		return comp.Left, bson.M{operator: v}, nil
	}

	if comp.Right[0] == '\'' {
		return comp.Left, bson.M{operator: comp.Right[1 : len(comp.Right)-1]}, nil
	}

	return comp.Left, bson.M{operator: comp.Right}, nil
}

func (stmt *Statement) ToMongoFind() (bson.D, bson.D, error) {
	if stmt.IsAggregate() {
		return bson.D{}, bson.D{}, fmt.Errorf("invalid statement for find")
	}

	selection, err := stmt.GetFindSelect()
	if err != nil {
		return bson.D{}, bson.D{}, err
	}

	if len(stmt.Where) == 0 {
		return bson.D{}, selection, nil
	}

	if stmt.BooleanOp == "" {
		k, v, err := stmt.GetFindWhere(0)
		if err != nil {
			return bson.D{}, bson.D{}, err
		}
		return bson.D{{k, v}}, selection, nil
	}

	booleanOpStr := ""

	if strings.ToUpper(stmt.BooleanOp) == "AND" {
		booleanOpStr = "$and"
	} else if strings.ToUpper(stmt.BooleanOp) == "OR" {
		booleanOpStr = "$or"
	} else {
		return bson.D{}, bson.D{}, fmt.Errorf("invalid boolean operator in WHERE")
	}

	arr := []bson.M{}
	for i := range stmt.Where {
		k, v, err := stmt.GetFindWhere(i)
		if err != nil {
			return bson.D{}, bson.D{}, err
		}

		arr = append(arr, bson.M{k: v})
	}

	return bson.D{{booleanOpStr, arr}}, selection, nil
}
