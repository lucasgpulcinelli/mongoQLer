package oracleManager

import "database/sql"

// struct UniqueEntry represents an unique constraint in oracle, having a base
// table and columns that generate the unique constraint.
type UniqueEntry struct {
	Table   string
	Columns []string
}

// GetUniques obtains all the unique constraints in the database as an array.
func GetUniques(db *sql.DB) ([]UniqueEntry, error) {

	// get all tables and columns that are of constraint type unique, ordering
	// by table to make sure when we change the table we are done with it
	s := `
    SELECT C.TABLE_NAME, CC.COLUMN_NAME
    FROM USER_CONSTRAINTS C
      JOIN USER_CONS_COLUMNS CC ON CC.CONSTRAINT_NAME = C.CONSTRAINT_NAME
    WHERE C.CONSTRAINT_TYPE = 'U'
    ORDER BY C.TABLE_NAME
  `

	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}

	result := []UniqueEntry{}
	tableNamePrev := ""

	// for each tuple
	for rows.Next() {
		tableName, columnName := "", ""

		// scan the table and column names
		err = rows.Scan(&tableName, &columnName)
		if err != nil {
			return nil, err
		}

		// if we are in the same table, the unique constraint is the same, so
		// append it to the existing entry
		if tableName == tableNamePrev {
			result[len(result)-1].Columns = append(
				result[len(result)-1].Columns, columnName,
			)
			continue
		}

		// if not, we are in a new table, so create a UniqueEntry for it
		result = append(
			result,
			UniqueEntry{Table: tableName, Columns: []string{columnName}},
		)

		tableNamePrev = tableName
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
