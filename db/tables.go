package db

import (
	"fmt"
	"database/sql"
//"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	projectSchema = `
create table project (
     id int primary key,
     name text,
     path text,
     mcid varchar(32)
)
`
	projectFilterSchema = `
create table project_filter (
     project_id 
)
`
)

func init() {
	fmt.Println("in db.tables init")
	db, err := sql.Open("sqlite3", "file:/tmp/sql.db?cached=shared&mode=rwc")
	if err != nil {
		fmt.Println("sql.Open err =", err)
	}
	_, err = db.Exec(projectSchema)
	if err != nil {
		fmt.Println("Exec err =", err)
	}
	db.Close()
}
