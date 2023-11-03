package sqlparser

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func (stmt *Statement) ToMongoAggregate() (mongo.Pipeline, error) {
	if !stmt.IsAggregate() {
		return mongo.Pipeline{}, fmt.Errorf("invalid statement for aggregation")
	}

	return mongo.Pipeline{}, nil
}
