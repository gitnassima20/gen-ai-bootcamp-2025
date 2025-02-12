package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

// TODO: Implement database connection and methods
