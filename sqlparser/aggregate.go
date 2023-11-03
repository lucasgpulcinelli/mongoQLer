package sqlparser

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
		k := ""

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
		}

		if k != "$count" {
			v := "$" + col.Name
			result = append(result, bson.E{col.Name, bson.D{{k, v}}})
		}
	}

	return result, nil
}

func (stmt *Statement) GetLookup() (bson.D, error) {
	if stmt.JoinTable == "" {
		return bson.D{}, nil
	}

	return bson.D{
		{"from", stmt.JoinTable},
		{"localField", stmt.JoinFromAttr},
		{"foreignField", stmt.JoinToAttr},
		{"as", fmt.Sprintf("%s_lookup", stmt.JoinTable)},
	}, nil
}

func (stmt *Statement) ToMongoAggregate() (mongo.Pipeline, error) {
	if !stmt.IsAggregate() {
		return mongo.Pipeline{}, fmt.Errorf("invalid statement for aggregation")
	}

	result := mongo.Pipeline{}

	where, err := stmt.GetFullWhere()
	if err != nil {
		return mongo.Pipeline{}, err
	}

	result = append(result, bson.D{{"$match", where}})

	join, err := stmt.GetLookup()
	if err != nil {
		return mongo.Pipeline{}, err
	}
	if len(join) != 0 {
		result = append(result, bson.D{{"$lookup", join}})
	}

	group, err := stmt.GetGroup()
	if err != nil {
		return mongo.Pipeline{}, err
	}

	if len(group) != 0 {
		result = append(result, bson.D{{"$group", group}})
	}

	return result, nil
}
