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

	// get all constraint names, tables and columns that are of constraint type
	// unique, ordering by constraint name to make it easier to differentiate
	// old from new UniqueEntries.
	s := `
    SELECT C.CONSTRAINT_NAME, C.TABLE_NAME, CC.COLUMN_NAME
    FROM USER_CONSTRAINTS C
      JOIN USER_CONS_COLUMNS CC ON CC.CONSTRAINT_NAME = C.CONSTRAINT_NAME
    WHERE C.CONSTRAINT_TYPE = 'U'
    ORDER BY C.CONSTRAINT_NAME
  `

	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}

	result := []UniqueEntry{}
	constraintNamePrev := ""

	// for each tuple
	for rows.Next() {
		var constraintName, tableName, columnName string

		// scan the constraint, table and column names
		err = rows.Scan(&constraintName, &tableName, &columnName)
		if err != nil {
			return nil, err
		}

		// if we are in the same constraint, the unique constraint is the same, so
		// append it to the existing entry
		if constraintName == constraintNamePrev {
			result[len(result)-1].Columns = append(
				result[len(result)-1].Columns, columnName,
			)
			continue
		}

		// if not, we are in a new unique, so create a UniqueEntry for it
		result = append(
			result,
			UniqueEntry{Table: tableName, Columns: []string{columnName}},
		)

		constraintNamePrev = constraintName
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
