package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"github.com/lucasgpulcinelli/mongoQLer/sqlparser"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	checkGeneratorButton *widget.Button
	mongoCGEntry         *widget.Entry
)

func checkGeneratorButtonFunc() {
	checks, err := oracleManager.GetChecks(oracleConn)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	s := ""
	for _, check := range checks {
		be, err := sqlparser.ParseBoolExpr(check.Check)
		if err != nil {
			errorPopUp(err, mainWindow.Canvas())
			return
		}

		bs, err := be.GetBson([]string{check.Table})
		if err != nil {
			errorPopUp(err, mainWindow.Canvas())
			return
		}

		bscomplete := bson.D{{"collMod", check.Table},
			{"validator", bs}, {"validationAction", "error"},
		}

		s += fmt.Sprintf("db.runCommand(%s)\n\n", bsonToString(bscomplete))
	}

	mongoCGEntry.SetText(s)
}

func newCheckGenerator() fyne.CanvasObject {

	checkGeneratorButton = widget.NewButton("generate",
		checkGeneratorButtonFunc,
	)

	mongoCGEntry = widget.NewMultiLineEntry()
	mongoCGEntry.SetPlaceHolder("click generate to get the validators")

	l := widget.NewLabel(
		"generate mongoDB validators from an oracle SQL connection",
	)

	return container.NewBorder(
		container.NewCenter(l), checkGeneratorButton, nil, nil, mongoCGEntry,
	)
}
