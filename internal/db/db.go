package db

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
)

var Conn *sql.DB

func Init() error {
	Conn, err := sql.Open("sqlite", "mist.db")
	if err != nil {
		return err
	}
	Conn.SetMaxIdleConns(1)
	Conn.SetMaxOpenConns(1)
	return dbMigrate(Conn)
}
