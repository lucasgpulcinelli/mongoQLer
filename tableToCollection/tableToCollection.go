package tableToCollection

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"go.mongodb.org/mongo-driver/bson"
)

// GetCollection generates, based on an oracle database connection, a series of
// documents with the data from a table or query. It will embed all references
// in the embeds array with the whole document related to the reference.
// This function takes in a database connection, an open query with data to be
// turned into the document, the table related to the query, and the list of
// references that should be embed.
//
// This function does recursion, and does not check if the embed references
// loop back to an initial table in order to stop infinite recursion.
func GetCollection(
	db *sql.DB, rows *sql.Rows, table string, embeds []oracleManager.Reference,
) ([]bson.D, error) {

	result := []bson.D{}

	// get column names
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// for each tuple
	for rows.Next() {
		// get its values (type does not matter)
		vs := make([]any, len(cols))
		for i := range cols {
			vs[i] = new(any)
		}

		err = rows.Scan(vs...)
		if err != nil {
			return nil, err
		}

		// and create the document with these values
		doc, err := writeDocument(db, table, embeds, cols, vs)
		if err != nil {
			return nil, err
		}

		// appending the document to the result array
		result = append(result, doc)
	}

	return result, nil
}

// writeDocument creates a mongoDB document from a list of columns and values
// from a certain table, embedding the references in the embeds array using the
// database connection provided.
// All primary keys in the tuple will be converted to a sub document in the
// "_id" field.
func writeDocument(
	db *sql.DB, table string, embeds []oracleManager.Reference, cols []string,
	vs []any,
) (bson.D, error) {

	// the final document
	doc := bson.D{}
	// the primary key sub document
	pks := bson.D{}

	// map to determine if the column provided is a reference
	refCols := map[string]bool{}

	// for each embedding reference
	for _, embed := range embeds {
		// if the reference does not relate to us, ignore it
		if embed.TableReferencer != table {
			continue
		}

		// if it does

		// for each column that generate the reference, get the reference values
		vsRef := []any{}
		for _, colRef := range embed.ColumnReferencer {
			refCols[colRef] = true

			for i, col := range cols {
				if col == colRef {
					// note that this wierd pointer conversion needs to happen because
					// to read data using rows.Scan we need an array of pointers and
					// to query data with db.Query using parameters we need an array of
					// values.

					vsRef = append(vsRef, *(vs[i].(*any)))
					break
				}
			}
		}

		// generate the subdocument with the reference value
		subDoc, err := embedReference(db, embeds, embed, vsRef)
		if err != nil {
			return nil, err
		}

		// and add it to our own document, generating the embedding
		doc = append(doc, bson.E{Key: embed.ConstraintName, Value: subDoc})
	}

	// for each column
	for i, name := range cols {
		if keyManager.IsPk([]string{table}, name) {
			// if it is a primary key, add it to the _id document
			pks = append(pks, bson.E{Key: name, Value: vs[i]})
		} else if refCols[name] != true {
			// if it is not a reference (in that case the embedding is already done)
			// add it to the main document
			doc = append(doc, bson.E{Key: name, Value: vs[i]})
		}
	}

	// if the primary key is empty, mongoDB should create an objectID for us, if
	// not, we should use our primary key
	if len(pks) != 0 {
		doc = append(doc, bson.E{Key: "_id", Value: pks})
	}

	return doc, nil
}

// embedReference creates a document containing matching data from a reference
// with certain tuple values
func embedReference(
	db *sql.DB, embeds []oracleManager.Reference, ref oracleManager.Reference,
	vs []any,
) (bson.D, error) {

	// if any value in the tuple is null, foreign key references do not apply,
	// so return a null reference
	for _, v := range vs {
		if v == nil {
			return nil, nil
		}
	}

	// create the query for getting the matching reference
	s := "SELECT * FROM " + ref.TableReferenced + " WHERE "

	for i, colName := range ref.ColumnReferenced {
		s += colName + " = :" + strconv.FormatInt(int64(i), 10) + " AND "
	}

	s = s[:len(s)-5]

	// execute the query with the tuple values
	rows, err := db.Query(s, vs...)
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

	// prepare the output tuple to receive data
	vsOut := make([]any, len(columns))
	for i := range columns {
		vsOut[i] = new(any)
	}

	// scan all columns
	err = rows.Scan(vsOut...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return nil, fmt.Errorf("embed query returned more than one tuple")
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	// and return a new document with the data from these tuples. If there are
	// any references that should be embed, this function will call
	// embedReference back, making the root of the embedding recursion.
	return writeDocument(
		db, ref.TableReferenced, embeds, columns, vsOut,
	)
}
