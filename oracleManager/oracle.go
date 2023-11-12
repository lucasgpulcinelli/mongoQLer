package oracleManager

import (
	"database/sql"
	"fmt"

	_ "github.com/sijms/go-ora/v2"
)

var (
	tableColumns map[string][]string
)

// Login logs in to a oracle database, returning the connection.
func Login(url, user, password string) (*sql.DB, error) {
	conn, err := sql.Open(
		"oracle",
		fmt.Sprintf("oracle://%s:%s@%s", user, password, url),
	)

	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	// initalise the table to columns map
	err = initTableColumnsMap(conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initTableColumnsMap(db *sql.DB) error {
	rows, err := db.Query("SELECT TABLE_NAME, COLUMN_NAME FROM USER_TAB_COLUMNS")
	if err != nil {
		return err
	}

	tableColumns = map[string][]string{}

	for rows.Next() {
		var table, column string

		err = rows.Scan(&table, &column)
		if err != nil {
			return err
		}

		tableColumns[table] = append(tableColumns[table], column)
	}

	return rows.Close()
}

func TableContainsColumn(table, columnToCheck string) bool {
	columns, ok := tableColumns[table]
	if !ok {
		return false
	}

	for _, col := range columns {
		if col == columnToCheck {
			return true
		}
	}

	return false
}

// GetPrimaryKeys gets all the primary keys in the databse in a map form taking
// the table name and returning all the column names that are primary keys.
func GetPrimaryKeys(db *sql.DB) (map[string][]string, error) {
	pks := map[string][]string{}

	s := `
  SELECT TABLE_NAME, COLUMN_NAME
  FROM USER_CONS_COLUMNS
    NATURAL JOIN USER_CONSTRAINTS
  WHERE CONSTRAINT_TYPE = 'P'
  `

	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t, c string

		err = rows.Scan(&t, &c)
		if err != nil {
			return nil, err
		}

		pks[t] = append(pks[t], c)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return pks, nil
}
