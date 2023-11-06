package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	mainWindow fyne.Window
)

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

func NewMainWindow(a fyne.App) fyne.Window {
	mainWindow = a.NewWindow("Oracle to Mongo Translator")

	tabs := container.NewAppTabs(
		container.NewTabItem("Table to Collection", widget.NewLabel("Hello, World!")),
		container.NewTabItem("Query to Collection", widget.NewLabel("Hello, World!")),
		container.NewTabItem("Indicies Generator", newIndiciesGenerator()),
		container.NewTabItem("Validator Generator", newCheckGenerator()),
		container.NewTabItem("Query to Find or Aggregate", newFindAggregate()),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	mainWindow.SetContent(tabs)
	mainWindow.Resize(fyne.NewSize(900, 500))

	return mainWindow
}
