package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	indiciesGeneratorButton *widget.Button
	mongoIGEntry            *widget.Entry
)

func indiciesGeneratorButtonFunc() {
	uniques, err := oracleManager.GetUniques(oracleConn)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	s := ""
	for _, un := range uniques {
		bs := bson.D{}
		for _, col := range un.Columns {
			bs = append(bs, bson.E{col, 1})
		}

		s += fmt.Sprintf("db.%s.createIndex(%s, {unique: true})\n\n",
			un.Table, bsonToString(bs),
		)
	}

	mongoIGEntry.SetText(s)
}

func newIndiciesGenerator() fyne.CanvasObject {

	indiciesGeneratorButton = widget.NewButton("generate",
		indiciesGeneratorButtonFunc,
	)

	mongoIGEntry = widget.NewMultiLineEntry()
	mongoIGEntry.SetPlaceHolder("click generate to get the indicies")

	l := widget.NewLabel(
		"generate mongoDB indicies from an oracle SQL connection",
	)

	return container.NewBorder(
		container.NewCenter(l), indiciesGeneratorButton, nil, nil, mongoIGEntry,
	)
}
