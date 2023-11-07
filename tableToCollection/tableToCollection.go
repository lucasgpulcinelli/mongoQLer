package tableToCollection

import (
	"database/sql"

	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCollection(
	db *sql.DB, table string, embed []oracleManager.Reference,
) ([]bson.D, error) {

	s := "SELECT * FROM " + table
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}

	result := []bson.D{}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		doc := bson.D{}
		vs := make([]any, len(cols))
		for i := range cols {
			vs[i] = new(any)
		}

		rows.Scan(vs...)

		if err != nil {
			return nil, err
		}

		pks := bson.D{}
		for i, name := range cols {
			if keyManager.IsPk([]string{table}, name) {
				pks = append(pks, bson.E{Key: name, Value: vs[i]})
			} else {
				doc = append(doc, bson.E{Key: name, Value: vs[i]})
			}
		}

		doc = append(doc, bson.E{Key: "_id", Value: pks})

		result = append(result, doc)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
