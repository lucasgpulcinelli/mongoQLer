package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"github.com/lucasgpulcinelli/mongoQLer/tableToCollection"
)

var (
	tableCollectionButton *widget.Button
	mongoTCEntry          *widget.Entry
	tcSelection           *widget.Select
	embedSelections       *fyne.Container
	referencesNow         []oracleManager.Reference
)

func initReferences(refs []oracleManager.Reference) {
	i := 0
	var constraint string
	for _, ref := range refs {
		if constraint == ref.ConstraintName {
			continue
		}

		constraint = ref.ConstraintName
		// and set them as checkbox options in one of the tabs
		j := i
		embedSelections.Add(container.NewHBox(
			widget.NewCheck("", func(set bool) {
				if !set {
					return
				}

				cont := embedSelections.Objects[j].(*fyne.Container)
				ck := cont.Objects[1].(*widget.Check)
				ck.SetChecked(false)
			}),
			widget.NewCheck(constraint, func(set bool) {
				if !set {
					return
				}

				cont := embedSelections.Objects[j].(*fyne.Container)
				ck := cont.Objects[0].(*widget.Check)
				ck.SetChecked(false)
			}),
		))
		i++
	}

}

// tableCollectionButtonFunc executes the table to collection button
// functionality.
func tableCollectionButtonFunc() {
	if referencesNow == nil || tcSelection.Selected == "" {
		return
	}

	// first, see wich references should be embedded based on the check boxes
	embedToRefs := []oracleManager.Reference{}
	embedFromRefs := []oracleManager.Reference{}
	for i, ref := range referencesNow {
		cont := embedSelections.Objects[i].(*fyne.Container)

		ckFrom := cont.Objects[0].(*widget.Check)
		ckTo := cont.Objects[1].(*widget.Check)

		if ckFrom.Checked {
			embedFromRefs = append(embedFromRefs, ref)
		} else if ckTo.Checked {
			embedToRefs = append(embedToRefs, ref)
		}
	}

	// execute the main query to pass to the GetCollection function
	rows, err := oracleConn.Query("SELECT * FROM " + tcSelection.Selected)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	// get all documents to be added to the new collection
	docs, err := tableToCollection.GetCollection(oracleConn, rows,
		tcSelection.Selected, embedToRefs, embedFromRefs,
	)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	// close the query, we are done by now
	err = rows.Close()
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	// format and print every new document to the mongo output text entry
	s := fmt.Sprintf("db.%s.insertMany([\n", tcSelection.Selected)
	i := 0
	for _, doc := range docs {
		s += bsonToString(doc) + ",\n"
		i++
	}
	if i > 0 {
		s = s[:len(s)-2] + "\n"
	}
	s += "])\n"

	mongoTCEntry.SetText(s)
}

// newTableToCollection creates the main object UI that converts an oracle
// table to a mongoDB collection. It takes as input a table and all references
// that should be embedded rather than linked, and returns in the mongoTCEntry
// a mongosh command to insert a new collection with all the data gathered from
// oracle.
func newTableToCollection() fyne.CanvasObject {
	l := widget.NewLabel("convert an oracle table to a mongoDB collection")
	l2 := widget.NewLabel("references to embed as arrays / objects")

	tableCollectionButton = widget.NewButton("convert",
		tableCollectionButtonFunc,
	)

	mongoTCEntry = widget.NewMultiLineEntry()
	mongoTCEntry.SetPlaceHolder("click convert to get your collection")

	// initialise the table selection as empty, the login flow will add the
	// tables after the connection is created.
	tcSelection = widget.NewSelect(
		[]string{},
		func(_ string) {},
	)

	embedSelections = container.NewVBox()

	return container.NewBorder(
		container.NewCenter(l),
		tableCollectionButton,
		nil,
		nil,
		container.NewHSplit(
			container.NewBorder(tcSelection, nil, nil, nil,
				container.NewBorder(container.NewCenter(l2), nil, nil, nil,
					container.NewVScroll(container.NewCenter(embedSelections)),
				),
			),
			mongoTCEntry,
		),
	)
}
