package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	mainWindow fyne.Window // the main ui window
)

var helpTC string = `
In this tab you can transform oracle tables to mongoDB collections.

You can select the table you want to convert via a selection box, as well as
all the foreign key references that should be embedded in your collection via
checkboxes.

The first checkbox means that the collection that is referenced will contain an
array of all the documents that reference it (for instance, marking this
option for an N-1 relation between city and state will make it so the state
contains an array of all their cities).
On the other hand the second checkbox means that the collection that creates
the reference will contain the referenced document as a whole (in the same
example, it would make it so each city of every state contains a copy of all of
the state's data).

Also, embedded documents can contain references to be embedded in them, the
program uses recursion into each document to embed everything as specified.

You can use each reference in one of the two forms described, but not both,
since it would create an infinite recursion.
`

var helpQC string = `
In this tab you can convert oracle query results into mongoDB collections.

You can write any oracle valid SQL query statement and the result is
transformed to a collection with the key value pairs as the SQL code describes
them. The _id of all tables created this way is a simple mongoDB ObjectID.

You can specify the output collection name in the box near the bottom of the
window.
`

var helpIG string = `
In this tab you can generate all the indicies for mongoDB collections based on
your oracle tables.

This tab generates a new index for each UNIQUE contraint in oracle (not primary
keys, since the _id already has an index in mongoDB).
`

var helpVG string = `
In this tab you can generate all the validtaors for mongoDB collections based
on your oracle tables.

This tab generates a new validator for each table / collection, considering
primary key constraints as well as check constraints (the first considering
that the primary key cannot be null, and including in the check constraints
all NOT NULLs as well as custom CHECKS).

The check parsing is done using the same SQL parser as in the last
functionality WHERE statement, so the same syntax restrictions for there apply
as well for all checks.

Note that, because the check parser has no context outside of the text itself,
it does differentiate between 1 and '1', for instance, even if the type for the
original table is a CHAR (which would make both the same in oracle).
`

var helpQFA string = `
In this tab you can convert an SQL query to a mongoDB find or aggregate.

The selection is in the form "SELECT A, B, C, ..." for simple column selection
or "SELECT GROUP_F_A(A), GROUP_F_B(B), ..." for group functions, where
GROUP_F_A or GROUP_F_B in this case can be any of SUM, AVG, COUNT (with support
for COUNT(*)), STDDEV, MIN and MAX. The output name for columns is the same as
it would be in the original document, and for group functions it is the name of
the column being selected (therefore different group functions in the same
column need to be done in different queries, as there is no alias support).
COUNT(*) is avaliable in the "count" attribute.

Only one table is avaliable using the FROM keyword.
For JOIN, only one inner non natual join is supported, with only one condition
as in "JOIN TABLE ON A = B", where A is an attribute from TABLE and B is an
attribute in the table defined in the FROM part of the query.

For WHERE conditions, simple comparisions using < > <= >= <> = as well as IN,
NOT IN, IS NULL and IS NOT NULL are supported. For combining these
comparisions, AND and OR can be used, but parenthesis must be used to mix them
(for instance, "A = B AND B = C AND C = D" is valid, but not
"A = B AND B = C OR C = D", for that we would need to use
"(A = B AND B = C) OR C = D").

There is no support for GROUP BY or ORDER BY.
`

// errorPopUp shows an error to a fyne canvas as a popup.
func errorPopUp(err error, c fyne.Canvas) {
	content := container.NewVBox(widget.NewLabel(fmt.Sprintf("error: %v", err)))

	pop := widget.NewModalPopUp(
		content,
		c,
	)

	content.Add(widget.NewButton("ok", func() {
		pop.Hide()
	}))

	pop.Show()
}

// NewMainWindow creates the main application window.
func NewMainWindow(a fyne.App) {
	mainWindow = a.NewWindow("Oracle to Mongo Translator")
	mainWindow.Resize(fyne.NewSize(900, 501))

	// the tabbed panes for each part of the application
	tabs := container.NewAppTabs(
		container.NewTabItem("Table to Collection", newTableToCollection()),
		container.NewTabItem("Query to Collection", newQueryToCollection()),
		container.NewTabItem("Indicies Generator", newIndiciesGenerator()),
		container.NewTabItem("Validator Generator", newCheckGenerator()),
		container.NewTabItem("Query to Find or Aggregate", newFindAggregate()),
	)

	mainWindow.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("help",
		fyne.NewMenuItem("about this tab", func() {
			var help string
			switch tabs.SelectedIndex() {
			default:
				return
			case 0:
				help = helpTC
			case 1:
				help = helpQC
			case 2:
				help = helpIG
			case 3:
				help = helpVG
			case 4:
				help = helpQFA
			}

			content := container.NewVBox(widget.NewLabel(help))

			pop := widget.NewModalPopUp(
				content,
				mainWindow.Canvas(),
			)

			content.Add(widget.NewButton("ok", func() {
				pop.Hide()
			}))

			pop.Show()
		})),
	))

	tabs.SetTabLocation(container.TabLocationLeading)

	go func() {
		time.Sleep(500 * time.Millisecond)
		mainWindow.Resize(fyne.NewSize(900, 500))
	}()

	mainWindow.SetContent(tabs)
}
