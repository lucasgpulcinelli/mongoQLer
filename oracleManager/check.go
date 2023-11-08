package oracleManager

import "database/sql"

// struct CheckEntry represents a check constraint in an oracle database,
// having a table that it applies to and a check condition itself.
type CheckEntry struct {
	Table string
	Check string
}

// GetChecks obtains all CheckEntries from a connection. It concatenates all
// the check conditions such that there is only one per table.
func GetChecks(db *sql.DB) ([]CheckEntry, error) {
	// obtain the table and search condition for all constraint of type check
	s := `
    SELECT C.TABLE_NAME, C.SEARCH_CONDITION_VC
    FROM USER_CONSTRAINTS C
    WHERE C.CONSTRAINT_TYPE = 'C'
    ORDER BY C.TABLE_NAME
  `

	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}

	result := []CheckEntry{}
	tableNamePrev := ""

	// for each check
	for rows.Next() {
		tableName, searchCondition := "", ""

		// scan it
		err = rows.Scan(&tableName, &searchCondition)
		if err != nil {
			return nil, err
		}

		// if we are still in the same table, concatenate the previous check with
		// the new condition
		if tableName == tableNamePrev {
			result[len(result)-1].Check += " AND (" + searchCondition + ")"
			continue
		}

		// if we changed tables, append a new CheckEntry
		result = append(
			result,
			CheckEntry{Table: tableName, Check: "(" + searchCondition + ")"},
		)

		tableNamePrev = tableName
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	for i := range result {
		result[i].Check += ";"
	}

	return result, nil
}
