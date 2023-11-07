package oracleManager

import "database/sql"

type CheckEntry struct {
	Table string
	Check string
}

func GetChecks(db *sql.DB) ([]CheckEntry, error) {
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

	for rows.Next() {
		tableName, searchCondition := "", ""
		err = rows.Scan(&tableName, &searchCondition)
		if err != nil {
			return nil, err
		}

		if tableName == tableNamePrev {
			result[len(result)-1].Check += " AND (" + searchCondition + ")"
			continue
		}

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
