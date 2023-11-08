package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/lucasgpulcinelli/mongoQLer/sqlparser"
)

var (
	findAggregateButton *widget.Button
	sqlFAEntry          *widget.Entry
	mongoFAEntry        *widget.Entry
)

// bsonToString transforms a bson to a string and treats erros using the UI
// popup.
func bsonToString(a bson.D) string {
	bts, err := bson.MarshalExtJSON(a, false, false)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return "{}"
	}

	return string(bts)
}

// findAggregateButtonFunc executes the SQL to find or aggregate functionality.
func findAggregateButtonFunc() {

	// first, parse the SQL
	stmt, err := sqlparser.Parse(sqlFAEntry.Text)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	if stmt.IsAggregate() {
		// get the aggregation from the statement
		mongoResult, err := stmt.ToMongoAggregate()
		if err != nil {
			errorPopUp(err, mainWindow.Canvas())
			return
		}

		// and format it for the final text output
		out := "[\n"
		for _, bs := range mongoResult {
			out += bsonToString(bs) + ",\n"
		}

		if len(mongoResult) > 1 {
			out = out[:len(out)-2]
		}

		mongoFAEntry.SetText(
			fmt.Sprint("db.", stmt.FromTable, ".aggregate(", out, "\n])"),
		)
	} else {
		// if the statement is a find

		// get the find and selection from the statement
		find, selection, err := stmt.ToMongoFind()
		if err != nil {
			errorPopUp(err, mainWindow.Canvas())
			return
		}

		// and format it for the final text output
		findJson := bsonToString(find)
		selectionJson := bsonToString(selection)

		mongoFAEntry.SetText(
			fmt.Sprint("db.", stmt.FromTable, ".find(\n", findJson, ",\n",
				selectionJson, "\n)"),
		)
	}
}

// newFindAggregate generates the main SQL query to find or aggregate mongoDB
// query. It takes the SQL text and outputs the mongoDB in another text area.
func newFindAggregate() fyne.CanvasObject {
	findAggregateButton = widget.NewButton("convert", findAggregateButtonFunc)

	sqlFAEntry = widget.NewMultiLineEntry()
	sqlFAEntry.SetText("SELECT * FROM DUAL;")

	mongoFAEntry = widget.NewMultiLineEntry()
	mongoFAEntry.SetPlaceHolder("click convert to convert your query")

	return container.NewBorder(
		container.NewCenter(
			widget.NewLabel("Convert an SQL query to a mongoDB find or aggregate"),
		),
		findAggregateButton,
		nil,
		nil,
		container.NewHSplit(sqlFAEntry, mongoFAEntry),
	)
}
