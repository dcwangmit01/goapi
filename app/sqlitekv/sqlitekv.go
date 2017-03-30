package sqlitekv

import (
	"database/sql"
	
	sqlkv "github.com/laurent22/go-sqlkv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbName = "app.db"
	tableName = "keyval"

	SqlKV *sqlkv.SqlKv
)

func Init() {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	
	// Will create table if not exists
	SqlKV = sqlkv.New(db, tableName)
}
