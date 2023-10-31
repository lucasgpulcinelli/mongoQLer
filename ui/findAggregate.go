package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/lucasgpulcinelli/mongoQLer/sqlparser"
)

var (
	findAggregateButton *widget.Button
	sqlEntry            *widget.Entry
	mongoEntry          *widget.Entry
)

func findAggregateButtonFunc() {
	stmt, err := sqlparser.Parse(sqlEntry.Text)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	mongoEntry.SetText(stmt.ToMongoFind())
}

func newFindAggregate() fyne.CanvasObject {
	findAggregateButton = widget.NewButton("convert", findAggregateButtonFunc)
	sqlEntry = widget.NewMultiLineEntry()
	sqlEntry.SetText("SELECT * FROM DUAL;")

	mongoEntry = widget.NewMultiLineEntry()

	return container.NewBorder(
		container.NewCenter(
			widget.NewLabel("Convert an SQL query to a mongoDB find or aggregate"),
		),
		findAggregateButton,
		nil,
		nil,
		container.NewHSplit(sqlEntry, mongoEntry),
	)
}
