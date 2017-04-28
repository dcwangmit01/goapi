package sqlitekv

import (
	"database/sql"

	sqlkv "github.com/laurent22/go-sqlkv"
	_ "github.com/mattn/go-sqlite3"
)

const (
	tableName = "keyval"
)

func New(dbPath string) *sqlkv.SqlKv {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	// Will create table if not exists
	return sqlkv.New(db, tableName)
}
