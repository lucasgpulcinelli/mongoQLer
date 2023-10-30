package main

import (
	"context"

	"fyne.io/fyne/v2/app"
	"github.com/lucasgpulcinelli/mongoQLer/ui"
)

func main() {
	a := app.New()
	wMain := ui.NewMainWindow(a)
	wLogin := ui.NewLoginWindow(a, wMain)

	wLogin.Show()
	a.Run()

	orConn, monConn := ui.GetConnections()
	if orConn != nil {
		err := orConn.Close()
		if err != nil {
			panic(err)
		}
	}
	if monConn != nil {
		err := monConn.Client().Disconnect(context.TODO())
		if err != nil {
			panic(err)
		}
	}
}
