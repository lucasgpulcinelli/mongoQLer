package ui

import (
	"context"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lucasgpulcinelli/mongoQLer/mongoManager"
	"github.com/lucasgpulcinelli/mongoQLer/oracleManager"

	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	oracleConn *sql.DB
	mongoConn  *mongo.Database
)

func GetConnections() (*sql.DB, *mongo.Database) {
	return oracleConn, mongoConn
}

func NewLoginWindow(a fyne.App, wMain fyne.Window) fyne.Window {
	w := a.NewWindow("Login to Oracle and MongoDB")
	w.Resize(fyne.NewSize(500, 500))

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

	mongoURL := widget.NewEntry()
	s, ok = os.LookupEnv("MONGO_URL")
	if !ok {
		s = "URL"
	}
	mongoURL.SetText(s)

	mongoDBName := widget.NewEntry()
	s, ok = os.LookupEnv("MONGO_DB_NAME")
	if !ok {
		s = "DB Name"
	}
	mongoDBName.SetText(s)

	mongoUser := widget.NewEntry()
	s, ok = os.LookupEnv("MONGO_USER")
	if !ok {
		s = "User"
	}
	mongoUser.SetText(s)

	mongoPass := widget.NewPasswordEntry()
	s, ok = os.LookupEnv("MONGO_PASSWORD")
	if !ok {
		s = "Password"
	}
	mongoPass.SetText(s)

	b := widget.NewButton("login", func() {
		var err error

		if oracleConn != nil {
			oracleConn.Close()
		}

		oracleConn, err = oracleManager.Login(oracleURL.Text, oracleUser.Text,
			oraclePass.Text,
		)
		if err != nil {
			errorPopUp(err, w.Canvas())
			return
		}

		if mongoConn != nil {
			mongoConn.Client().Disconnect(context.TODO())
		}

		mongoConn, err = mongoManager.Login(mongoURL.Text, mongoDBName.Text,
			mongoUser.Text, mongoPass.Text,
		)
		if err != nil {
			errorPopUp(err, w.Canvas())
			return
		}

		w.Close()
		wMain.Show()
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

		widget.NewSeparator(),

		container.NewBorder(
			nil, nil, widget.NewLabel("MongoDB URL"), nil, mongoURL,
		),
		container.NewBorder(
			nil, nil, widget.NewLabel("MongoDB Database Name"), nil, mongoDBName,
		),
		container.NewBorder(
			nil, nil, widget.NewLabel("MongoDB User"), nil, mongoUser,
		),
		container.NewBorder(
			nil, nil, widget.NewLabel("MongoDB Password"), nil, mongoPass,
		),
		b,
	)

	w.SetContent(content)

	return w
}
