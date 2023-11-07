package tableToCollection

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCollection(
	db *sql.DB, rows *sql.Rows, table string, embeds []oracleManager.Reference,
) ([]bson.D, error) {

	result := []bson.D{}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		vs := make([]any, len(cols))
		for i := range cols {
			vs[i] = new(any)
		}

		rows.Scan(vs...)

		if err != nil {
			return nil, err
		}

		doc, err := writeDocument(db, table, embeds, cols, vs)
		if err != nil {
			return nil, err
		}

		result = append(result, doc)
	}

	return result, nil
}

func writeDocument(
	db *sql.DB, table string, embeds []oracleManager.Reference, cols []string,
	vs []any,
) (bson.D, error) {

	doc := bson.D{}
	pks := bson.D{}

	refCols := map[string]bool{}

	for _, embed := range embeds {
		if embed.TableReferencer != table {
			continue
		}

		vsRef := []any{}
		for _, colRef := range embed.ColumnReferencer {
			refCols[colRef] = true

			for i, col := range cols {
				if col == colRef {
					vsRef = append(vsRef, vs[i])
					break
				}
			}
		}

		subDoc, err := embedReference(db, embeds, embed, vsRef)
		if err != nil {
			return nil, err
		}

		doc = append(doc, bson.E{Key: embed.ConstraintName, Value: subDoc})
	}

	for i, name := range cols {
		if keyManager.IsPk([]string{table}, name) {
			pks = append(pks, bson.E{Key: name, Value: vs[i]})
		} else if refCols[name] != true {
			doc = append(doc, bson.E{Key: name, Value: vs[i]})
		}
	}

	if len(pks) != 0 {
		doc = append(doc, bson.E{Key: "_id", Value: pks})
	}

	return doc, nil
}

func embedReference(
	db *sql.DB, embeds []oracleManager.Reference, ref oracleManager.Reference,
	vs []any,
) (bson.D, error) {
	for _, v := range vs {
		vv := v.(*any)
		if *vv == nil {
			return nil, nil
		}
	}

	s := "SELECT * FROM " + ref.TableReferenced + " WHERE "

	for i, colName := range ref.ColumnReferenced {
		s += colName + " = :" + strconv.FormatInt(int64(i), 10) + " AND "
	}

	s = s[:len(s)-5]

	vp := []any{}
	for _, v := range vs {
		vp = append(vp, *(v.(*any)))
	}

	rows, err := db.Query(s, vp...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("embed query returned no data")
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	vsOut := make([]any, len(columns))
	for i := range columns {
		vsOut[i] = new(any)
	}

	err = rows.Scan(vsOut...)
	if err != nil {
		return nil, err
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return writeDocument(
		db, ref.TableReferenced, embeds, columns, vsOut,
	)
}
