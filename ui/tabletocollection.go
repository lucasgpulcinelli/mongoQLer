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

func tableCollectionButtonFunc() {
	if referencesNow == nil || tcSelection.Selected == "" {
		return
	}

	embedRefs := []oracleManager.Reference{}
	for i, ref := range referencesNow {
		checkBox, ok := embedSelections.Objects[i].(*widget.Check)
		if !ok {
			errorPopUp(fmt.Errorf("embedSelections has wrong widget types"),
				mainWindow.Canvas(),
			)
			return
		}

		if !checkBox.Checked {
			continue
		}

		embedRefs = append(embedRefs, ref)
	}

	rows, err := oracleConn.Query("SELECT * FROM " + tcSelection.Selected)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	docs, err := tableToCollection.GetCollection(oracleConn, rows,
		tcSelection.Selected, embedRefs,
	)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	err = rows.Close()
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

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

func newTableToCollection() fyne.CanvasObject {
	l := widget.NewLabel("convert an oracle table to a mongoDB collection")
	l2 := widget.NewLabel("references to embed")

	tableCollectionButton = widget.NewButton("convert",
		tableCollectionButtonFunc,
	)

	mongoTCEntry = widget.NewMultiLineEntry()
	mongoTCEntry.SetPlaceHolder("click convert to get your collection")

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
					container.NewVScroll(embedSelections),
				),
			),
			mongoTCEntry,
		),
	)
}
