package oracleManager

import "database/sql"

type UniqueEntry struct {
	Table   string
	Columns []string
}

func GetUniques(db *sql.DB) ([]UniqueEntry, error) {
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

	for rows.Next() {
		tableName, columnName := "", ""
		
    err = rows.Scan(&tableName, &columnName)
    if err != nil {
      return nil, err
    }

		if tableName == tableNamePrev {
			result[len(result)-1].Columns = append(
				result[len(result)-1].Columns, columnName,
			)
			continue
		}

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
