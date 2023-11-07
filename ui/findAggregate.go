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

func bsonToString(a bson.D) string {
	bts, err := bson.MarshalExtJSON(a, false, false)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return "{}"
	}

	return string(bts)
}

func findAggregateButtonFunc() {
	stmt, err := sqlparser.Parse(sqlFAEntry.Text)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	if stmt.IsAggregate() {
		mongoResult, err := stmt.ToMongoAggregate()
		if err != nil {
			errorPopUp(err, mainWindow.Canvas())
			return
		}

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
		find, selection, err := stmt.ToMongoFind()
		if err != nil {
			errorPopUp(err, mainWindow.Canvas())
			return
		}

		findJson := bsonToString(find)
		selectionJson := bsonToString(selection)

		mongoFAEntry.SetText(
			fmt.Sprint("db.", stmt.FromTable, ".find(\n", findJson, ",\n",
				selectionJson, "\n)"),
		)
	}
}

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
