package sqlparser

import (
	"fmt"

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

func (stmt *Statement) ToMongoFind() (bson.D, bson.D, error) {
	if stmt.IsAggregate() {
		return bson.D{}, bson.D{}, fmt.Errorf("invalid statement for find")
	}

	selection, err := stmt.GetFindSelect()
	if err != nil {
		return bson.D{}, bson.D{}, err
	}

	where, err := stmt.Where.GetBson()
	if err != nil {
		return bson.D{}, bson.D{}, err
	}

	return where, selection, nil
}
