package db

import (
//"database/sql"
//"github.com/jmoiron/sqlx"
//"github.com/mattn/go-sqlite3"
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
