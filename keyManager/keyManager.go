package keyManager

import (
	"database/sql"

	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
)

var primaryKeys map[string][]string

func InitPrimaryKeys(db *sql.DB) (err error) {
	primaryKeys, err = oracleManager.GetPrimaryKeys(db)
	return err
}

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
