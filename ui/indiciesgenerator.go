package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	indiciesGeneratorButton *widget.Button
	mongoIGEntry            *widget.Entry
)

// indiciesGeneratorButton executes the indicies generator button
// functionality.
func indiciesGeneratorButtonFunc() {
	// get all the unique entries from oracle
	uniques, err := oracleManager.GetUniques(oracleConn)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	// for each unique, convert it to a document and add an createIndex for it
	s := ""
	for _, un := range uniques {
		bs := bson.D{}
		for _, col := range un.Columns {
			bs = append(bs, bson.E{
				Key:   keyManager.ToMongoId(un.Table, col),
				Value: 1,
			})
		}

		s += fmt.Sprintf("db.%s.createIndex(%s,\n{unique: true}\n)\n",
			un.Table, bsonToString(bs),
		)
	}

	mongoIGEntry.SetText(s)
}

// newIndiciesGenerator generates the main indicies generator UI. It generates
// all the indicies for a mongoDB database from UNIQUE constraints from oracle.
func newIndiciesGenerator() fyne.CanvasObject {

	indiciesGeneratorButton = widget.NewButton("generate",
		indiciesGeneratorButtonFunc,
	)

	mongoIGEntry = widget.NewMultiLineEntry()
	mongoIGEntry.SetPlaceHolder("click generate to get the indicies")

	l := widget.NewLabel(
		"Generate mongoDB indicies from an oracle SQL connection",
	)

	return container.NewBorder(
		container.NewCenter(l), indiciesGeneratorButton, nil, nil, mongoIGEntry,
	)
}
