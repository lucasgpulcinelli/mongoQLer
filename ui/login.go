package ui

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lucasgpulcinelli/mongoQLer/keyManager"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"

	"database/sql"
)

var (
	oracleConn *sql.DB
)

// GetConnections gets the oracle database connection. Before the login window
// closes the connection is nil.
func GetConnections() *sql.DB {
	return oracleConn
}

// NewLoginWindow creates the login window for connecting with the database.
// The main window will be activated once the login flow is complete.
// The application uses as initial text the values from some environment
// variables, such that if a .env file is created and sourced, it is easier to
// rerun the application multiple times.
func NewLoginWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Login to Oracle")
	w.Resize(fyne.NewSize(500, 300))

	oracleURL := widget.NewEntry()
	s, ok := os.LookupEnv("ORACLE_URL")
	if !ok {
		s = "URL"
	}
	oracleURL.SetText(s)

	oracleUser := widget.NewEntry()
	s, ok = os.LookupEnv("ORACLE_USER")
	if !ok {
		s = "User"
	}
	oracleUser.SetText(s)

	oraclePass := widget.NewPasswordEntry()
	s, ok = os.LookupEnv("ORACLE_PASSWORD")
	if !ok {
		s = "Password"
	}
	oraclePass.SetText(s)

	// the button that will execute the login functionality
	b := widget.NewButton("login", func() {
		var err error

		if oracleConn != nil {
			oracleConn.Close()
		}

		// log in to oracle
		oracleConn, err = oracleManager.Login(oracleURL.Text, oracleUser.Text,
			oraclePass.Text,
		)
		if err != nil {
			errorPopUp(err, w.Canvas())
			return
		}

		// initialise some oracle metadata:

		// get all foreign key references
		referencesNow, err = oracleManager.GetReferences(oracleConn)
		if err != nil {
			errorPopUp(err, w.Canvas())
			return
		}

		// get all primary keys and prepare the key manager
		err = keyManager.InitPrimaryKeys(oracleConn)
		if err != nil {
			errorPopUp(err, w.Canvas())
			return
		}

		// get all table names
		tables, err := oracleManager.GetTables(oracleConn)
		if err != nil {
			errorPopUp(err, w.Canvas())
			return
		}

		NewMainWindow(a)

		tcSelection.SetOptions(tables)
		initReferences(referencesNow)

		mainWindow.Show()
		w.Close()
	})

	content := container.NewVBox(
		container.NewBorder(
			nil, nil, widget.NewLabel("Oracle URL"), nil, oracleURL,
		),
		container.NewBorder(
			nil, nil, widget.NewLabel("Oracle User"), nil, oracleUser,
		),
		container.NewBorder(
			nil, nil, widget.NewLabel("Oracle Password"), nil, oraclePass,
		),
		b,
	)

	w.SetContent(content)

	return w
}
