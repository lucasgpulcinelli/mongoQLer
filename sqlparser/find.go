package sqlparser

import (
	"fmt"
	"strings"

	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"go.mongodb.org/mongo-driver/bson"
)

// GetSelect gets the bson representing the find key selection document, or the
// $project value in an aggregation pipeline.
func (stmt *Statement) GetSelect() (bson.D, error) {
	if len(stmt.SelectColumn) == 0 {
		return bson.D{}, nil
	}

	ret := bson.D{}

	hasKey := false
	for _, selection := range stmt.SelectColumn {
		var k string

		// if the column is in the joined table, we need to reference it as
		// table.column, because the lookup + unwind will make the reference to
		// this field as that
		if oracleManager.TableContainsColumn(stmt.JoinTable, selection.Name) &&
			selection.GroupFunction == "" {

			k = stmt.JoinTable + "." + keyManager.ToMongoId(
				stmt.JoinTable,
				selection.Name,
			)
		} else if selection.GroupFunction != "" {
			k = selection.Name
		} else {
			k = keyManager.ToMongoId(
				stmt.FromTable,
				selection.Name,
			)
		}

		// if we are using the _id from the FromTable, mark that we cannot omit
		// the _id field
		if len(k) > 3 && k[:3] == "_id" {
			hasKey = true
		}

		// the COUNT group funcion is special, we need to use count as the
		// selection key
		if k == "*" && strings.ToUpper(selection.GroupFunction) == "COUNT" {
			ret = append(ret, bson.E{Key: "count", Value: 1})
			continue
		}

		ret = append(ret, bson.E{Key: k, Value: 1})
	}

	// if we don't have any keys involved, mongoDB will assume we want them, so
	// explictly mark that we don't want the key
	if len(ret) != 0 && !hasKey {
		ret = append(ret, bson.E{Key: "_id", Value: 0})
	}

	return ret, nil
}

// ToMongoFind gets the bsons representing a find for a statement. The first
// document is the filter and the second is the key selection.
func (stmt *Statement) ToMongoFind() (bson.D, bson.D, error) {
	if stmt.IsAggregate() {
		return bson.D{}, bson.D{}, fmt.Errorf("invalid statement for find")
	}

	selection, err := stmt.GetSelect()
	if err != nil {
		return bson.D{}, bson.D{}, err
	}

	where, err := stmt.Where.GetBson(stmt.FromTable, stmt.JoinTable)
	if err != nil {
		return bson.D{}, bson.D{}, err
	}

	return where, selection, nil
}
