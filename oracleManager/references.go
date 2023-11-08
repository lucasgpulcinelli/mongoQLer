package oracleManager

import (
	"database/sql"
)

// struct Reference represents an oracle foreign key constraint from a table
// and columns to another.
type Reference struct {
	ConstraintName   string
	TableReferencer  string
	ColumnReferencer []string
	TableReferenced  string
	ColumnReferenced []string
}

// GetReferences gets all the foreign key references as an array.
func GetReferences(db *sql.DB) ([]Reference, error) {
	// get the constraint name, table and columns that generate the reference,
	// and table and columns that are referenced for each constraint of type
	// reference, ordered by constraint name to keep track of our current
	// reference easily.
	s := `
  SELECT C.CONSTRAINT_NAME, 
    CC.TABLE_NAME AS TABLE_REFERENCER, CC.COLUMN_NAME AS COLUMN_REFERENCER,
    RCC.TABLE_NAME AS TABLE_REFERENCED, RCC.COLUMN_NAME AS COLUMN_REFERENCED
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

	// for each reference
	ret := []Reference{}
	prevCn := ""
	for rows.Next() {
		var cn, tr, cr, td, cd string

		// scan it
		err = rows.Scan(&cn, &tr, &cr, &td, &cd)
		if err != nil {
			return nil, err
		}

		// if the constraint name changed, create a new entry in the array
		if prevCn != cn {
			prevCn = cn
			ret = append(ret, Reference{cn, tr, []string{cr}, td, []string{cd}})
			continue
		}

		// if not, append it to the previous reference: we are still in the same
		// composite foreign key.

		ret[len(ret)-1].ColumnReferencer = append(
			ret[len(ret)-1].ColumnReferencer, cr,
		)

		ret[len(ret)-1].ColumnReferenced = append(
			ret[len(ret)-1].ColumnReferenced, cd,
		)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// GetTables obtains all tables from an oracle connection and returns them in
// an array.
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
