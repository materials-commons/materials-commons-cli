package materials

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type fileDB struct {
	*leveldb.DB
}

func openFileDB(path string) (*fileDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &fileDB{DB: db}, nil
}
