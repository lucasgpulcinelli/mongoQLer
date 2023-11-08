package keyManager

import (
	"database/sql"

	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
)

// a map from a table to all the columns that are keys for that table
var primaryKeys map[string][]string

// InitPrimaryKeys reads all primary key definitions from all tables and adds
// them to the primaryKeys map, enabling the use of all other functions from
// the package.
func InitPrimaryKeys(db *sql.DB) (err error) {
	primaryKeys, err = oracleManager.GetPrimaryKeys(db)
	return err
}

// IsPk indicates if a column is a key in any of the tables provided
func IsPk(tables []string, column string) bool {
	for _, table := range tables {
		for _, pkc := range primaryKeys[table] {
			if pkc == column {
				return true
			}
		}
	}

	return false
}

// ToMongoId converts a column to a mongoDB key for a document, converting the
// column name to "_id.columnname" if the column is a primary key from oracle
// in one of the tables provided.
func ToMongoId(tables []string, column string) string {
	for _, table := range tables {
		for _, pkc := range primaryKeys[table] {
			if pkc == column {
				return "_id." + column
			}
		}
	}

	return column
}
