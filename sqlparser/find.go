package sqlparser

import (
	"fmt"
	"strings"

	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
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
		k := keyManager.ToMongoId(
			[]string{stmt.FromTable},
			selection.Name,
		)

		if strings.Contains(k, "_id") {
			hasKey = true
		}

		if k == "*" && strings.ToUpper(selection.GroupFunction) == "COUNT" {
			ret = append(ret, bson.E{"count", 1})
			continue
		}

		ret = append(ret, bson.E{k, 1})
	}

	if len(ret) != 0 && !hasKey {
		ret = append(ret, bson.E{"_id", 0})
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

	where, err := stmt.Where.GetBson([]string{stmt.FromTable, stmt.JoinTable})
	if err != nil {
		return bson.D{}, bson.D{}, err
	}

	return where, selection, nil
}
