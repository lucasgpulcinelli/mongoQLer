package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	w := a.NewWindow("Oracle to Mongo Translator")

	tabs := container.NewAppTabs(
		container.NewTabItem("Table to Collection", widget.NewLabel("Hello, World!")),
		container.NewTabItem("Query to Collection", widget.NewLabel("Hello, World!")),
		container.NewTabItem("Indicies Generator", widget.NewLabel("Hello, World!")),
		container.NewTabItem("Validator Generator", widget.NewLabel("Hello, World!")),
		container.NewTabItem("Query to Find and Aggregate", widget.NewLabel("Hello, World!")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(900, 500))

	return w
}
