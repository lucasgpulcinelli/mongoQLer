package sqlparser

import (
	"fmt"
	"strings"

	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetGroup gets the document paired with the $group operator for a mongDB
// aggregation for a Statement.
func (stmt *Statement) GetGroup() (bson.D, error) {
	hasGroup := false
	for _, col := range stmt.SelectColumn {
		if col.GroupFunction != "" {
			hasGroup = true
			break
		}
	}

	if !hasGroup {
		return bson.D{}, nil
	}

	result := bson.D{{"_id", nil}}

	for _, col := range stmt.SelectColumn {
		var k, v string

		switch strings.ToUpper(col.GroupFunction) {
		default:
			return bson.D{}, fmt.Errorf("invalid group function name")
		case "SUM":
			k = "$sum"
		case "MIN":
			k = "$min"
		case "MAX":
			k = "$max"
		case "AVG":
			k = "$avg"
		case "MEADIAN":
			k = "$median"
		case "COUNT":
			k = "$count"
			if col.Name != "*" {
				return bson.D{}, fmt.Errorf("COUNT is only supported as COUNT(*)")
			}
			result = append(result, bson.E{"count", bson.D{{k, bson.D{}}}})
			continue
		}

		if oracleManager.TableContainsColumn(stmt.JoinTable, col.Name) {
			v = "$" + stmt.JoinTable + "." +
				keyManager.ToMongoId(stmt.JoinTable, col.Name)
		} else {
			v = "$" + keyManager.ToMongoId(stmt.FromTable, col.Name)
		}

		result = append(result, bson.E{col.Name, bson.D{{k, v}}})
	}

	return result, nil
}

// GetLookup gets the document paired with the $lookup operator for a mongoDB
// aggregation for a Statement.
func (stmt *Statement) GetLookup() (bson.D, error) {
	if stmt.JoinTable == "" {
		return bson.D{}, nil
	}

	return bson.D{
		{"from", stmt.JoinTable},
		{"localField", keyManager.ToMongoId(stmt.FromTable, stmt.JoinFromAttr)},
		{"foreignField", keyManager.ToMongoId(stmt.JoinTable, stmt.JoinToAttr)},
		{"as", stmt.JoinTable},
	}, nil
}

// ToMongoAggregate gets the Pipeline representing an aggregation for a
// Statement.
func (stmt *Statement) ToMongoAggregate() (mongo.Pipeline, error) {
	if !stmt.IsAggregate() {
		return mongo.Pipeline{}, fmt.Errorf("invalid statement for aggregation")
	}

	result := mongo.Pipeline{}

	join, err := stmt.GetLookup()
	if err != nil {
		return mongo.Pipeline{}, err
	}
	if len(join) != 0 {
		result = append(result,
			bson.D{{"$lookup", join}},
			bson.D{{"$unwind", "$" + stmt.JoinTable}},
		)
	}

	where, err := stmt.Where.GetBson(stmt.FromTable, stmt.JoinTable)
	if err != nil {
		return mongo.Pipeline{}, err
	}

	result = append(result, bson.D{{"$match", where}})

	group, err := stmt.GetGroup()
	if err != nil {
		return mongo.Pipeline{}, err
	}

	if len(group) != 0 {
		result = append(result, bson.D{{"$group", group}})
	}

	selection, err := stmt.GetSelect()
	if err != nil {
		return mongo.Pipeline{}, err
	}

	if len(selection) != 0 {
		result = append(result, bson.D{{"$project", selection}})
	}

	return result, nil
}
