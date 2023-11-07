package sqlparser

import (
	"fmt"

	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
	"go.mongodb.org/mongo-driver/bson"
)

func (stmt *Statement) GetFindSelect() (bson.D, error) {
	if len(stmt.SelectColumn) == 0 {
		return bson.D{}, nil
	}

	ret := bson.D{}

	for _, selection := range stmt.SelectColumn {
		k := keyManager.ToMongoId(
			[]string{stmt.FromTable},
			selection.Name,
		)

		ret = append(ret, bson.E{k, 1})
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

	where, err := stmt.Where.GetBson([]string{stmt.FromTable, stmt.JoinTable})
	if err != nil {
		return bson.D{}, bson.D{}, err
	}

	return where, selection, nil
}
