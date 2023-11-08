package oracleManager

import (
	"database/sql"
	"fmt"

	_ "github.com/sijms/go-ora/v2"
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

	return conn, nil
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
