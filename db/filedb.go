package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type FileDB struct {
	*leveldb.DB
}

func OpenFileDB(path string) (*FileDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &FileDB{DB: db}, nil
}
