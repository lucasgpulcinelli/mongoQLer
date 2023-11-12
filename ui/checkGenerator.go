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

// checkGeneratorButtonFunc executes the check generation button functionality.
func checkGeneratorButtonFunc() {
	// get all checks from the oracle connection
	checks, err := oracleManager.GetChecks(oracleConn)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	s := ""

	// for each check
	for _, check := range checks {
		// parse it
		be, err := sqlparser.ParseBoolExpr(check.Check)
		if err != nil {
			errorPopUp(err, mainWindow.Canvas())
			return
		}

		// get the bson for that check
		bs, err := be.GetBson(check.Table, "")
		if err != nil {
			errorPopUp(err, mainWindow.Canvas())
			return
		}

		// generate the complete validator bson
		bscomplete := bson.D{{Key: "collMod", Value: check.Table},
			{Key: "validator", Value: bs}, {Key: "validationAction", Value: "error"},
		}

		// and concatenate it for the final text output
		s += fmt.Sprintf("db.runCommand(%s)\n", bsonToString(bscomplete))
	}

	mongoCGEntry.SetText(s)
}

// newCheckGenerator generates the main check generation UI. It generates all
// the validators for a mongoDB database from CHECK constraints from oracle.
func newCheckGenerator() fyne.CanvasObject {

	checkGeneratorButton = widget.NewButton("generate",
		checkGeneratorButtonFunc,
	)

	mongoCGEntry = widget.NewMultiLineEntry()
	mongoCGEntry.SetPlaceHolder("click generate to get the validators")

	l := widget.NewLabel(
		"Generate mongoDB validators from an oracle SQL connection",
	)

	return container.NewBorder(
		container.NewCenter(l), checkGeneratorButton, nil, nil, mongoCGEntry,
	)
}
