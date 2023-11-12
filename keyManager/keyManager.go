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

// IsPk indicates if a column is a key in the table provided
func IsPk(table string, column string) bool {
	for _, pkc := range primaryKeys[table] {
		if pkc == column {
			return true
		}
	}

	return false
}

// ToMongoId converts a column to a mongoDB key for a document, converting the
// column name to "_id.columnname" if the column is a primary key from oracle
// in the table provided.
func ToMongoId(table string, column string) string {
	for _, pkc := range primaryKeys[table] {
		if pkc == column {
			return "_id." + column
		}
	}

	return column
}
