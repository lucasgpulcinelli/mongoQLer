package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/lucasgpulcinelli/mongoQLer/ui"
)

func main() {
	a := app.New()
	wLogin := ui.NewLoginWindow(a)

	wLogin.Show()
	a.Run()

	// during the exit phase of the application, close the connection
	conn := ui.GetConnections()
	if conn != nil {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}
}
