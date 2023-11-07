package oracleManager

import (
	"database/sql"
)

type Reference struct {
	ConstraintName  string
	TableReferencer string
	TableReferenced string
}

func GetReferences(db *sql.DB) ([]Reference, error) {
	s := `
  SELECT C.CONSTRAINT_NAME, 
    CC.TABLE_NAME AS TABLE_REFERENCER, 
    RCC.TABLE_NAME AS TABLE_REFERENCED
  FROM USER_CONSTRAINTS C
    JOIN USER_CONS_COLUMNS CC ON
      CC.CONSTRAINT_NAME = C.CONSTRAINT_NAME
    JOIN USER_CONS_COLUMNS RCC ON
      C.R_CONSTRAINT_NAME = RCC.CONSTRAINT_NAME
      AND RCC.POSITION = CC.POSITION
  WHERE C.CONSTRAINT_TYPE = 'R'
  ORDER BY CONSTRAINT_NAME
  `

	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}

	ret := []Reference{}
	for rows.Next() {
		var cn, tr, td string

		err = rows.Scan(&cn, &tr, &td)
		if err != nil {
			return nil, err
		}

		ret = append(ret, Reference{cn, tr, td})
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func GetTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT TABLE_NAME FROM USER_TABLES")
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			return nil, err
		}

		ret = append(ret, s)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return ret, nil
}