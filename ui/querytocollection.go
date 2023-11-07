package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"
	"github.com/lucasgpulcinelli/mongoQLer/tableToCollection"
)

var (
	queryButton    *widget.Button
	sqlQCEntry     *widget.Entry
	mongoQCEntry   *widget.Entry
	queryNameEntry *widget.Entry
)

func queryButtonFunc() {
	if queryNameEntry.Text == "" {
		return
	}

	query := sqlQCEntry.Text
	if i := strings.Index(query, ";"); i != -1 {
		query = query[:i]
	}

	rows, err := oracleConn.Query(query)
	if err != nil {
		errorPopUp(err, mainWindow.Canvas())
		return
	}

	docs, err := tableToCollection.GetCollection(oracleConn, rows, "",
		[]oracleManager.Reference{},
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

	s := fmt.Sprintf("db.%s.insertMany([\n", queryNameEntry.Text)
	i := 0
	for _, doc := range docs {
		s += bsonToString(doc) + ",\n"
		i++
	}
	if i > 0 {
		s = s[:len(s)-2] + "\n"
	}
	s += "])\n"

	mongoQCEntry.SetText(s)
}

func newQueryToCollection() fyne.CanvasObject {
	queryButton = widget.NewButton("generate", queryButtonFunc)

	sqlQCEntry = widget.NewMultiLineEntry()
	sqlQCEntry.SetText("SELECT * FROM DUAL;")

	mongoQCEntry = widget.NewMultiLineEntry()
	mongoQCEntry.SetPlaceHolder("click generate to get your collection")

	queryNameEntry = widget.NewEntry()
	queryNameEntry.SetPlaceHolder("name of output collection")

	return container.NewBorder(
		container.NewCenter(
			widget.NewLabel("Convert an SQL query to a mongoDB collection"),
		),
		container.NewBorder(nil, nil, nil, queryButton, queryNameEntry),
		nil,
		nil,
		container.NewHSplit(sqlQCEntry, mongoQCEntry),
	)
}
