package db

import (
	"github.com/materials-commons/gohandy/file"
	_ "github.com/mattn/go-sqlite3"
)

func Exists(path string) bool {
	return file.Exists(path)
}
