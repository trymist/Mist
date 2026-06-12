package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB

func Init() error {
	conn, err := sql.Open("sqlite", "mist.db")
	Conn = conn
	if err != nil {
		return err
	}
	Conn.SetMaxIdleConns(1)
	Conn.SetMaxOpenConns(1)
	err = Conn.Ping()
	if err != nil {
		return err
	}
	return dbMigrate(Conn)
}
