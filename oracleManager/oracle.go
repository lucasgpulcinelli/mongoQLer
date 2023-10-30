package oracleManager

import (
	"database/sql"
	"fmt"

	_ "github.com/sijms/go-ora/v2"
)

func Login(url, user, password string) (*sql.DB, error) {
	conn, err := sql.Open(
		"oracle",
		fmt.Sprintf("oracle://%s:%s@%s", user, password, url),
	)

	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
