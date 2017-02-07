package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databasePath = "./sqlite.db"
)

var createSql = `
CREATE TABLE IF NOT EXISTS history (
    date VARCHAR(255) NOT NULL,
    dir  VARCHAR(255) NOT NULL,
    cmd  VARCHAR(255) NOT NULL
)`

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(createSql); err != nil {
		return nil, err
	}

	return db, nil
}
