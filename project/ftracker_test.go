package project

import (
	"testing"
	"github.com/materials-commons/materials"
	"github.com/syndtr/goleveldb/leveldb"
	"encoding/json"
	"fmt"
)

func TestWalkProject(t *testing.T) {
	if false {
		return
	}
	p := materials.Project{
		Name: "Data_MatComm2",
		Path: "/Users/gtarcea/Dropbox/transfers/materialscommons/Data_MatComm2",
		Status: "Unknown",
	}

	WalkProject(p)
}

func TestCreatedDb(t *testing.T) {
	db, _ := leveldb.OpenFile("/tmp/project.db", nil)
	iter := db.NewIterator(nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		var p projectFileInfo
		json.Unmarshal(value, &p)
		fmt.Println(key)
		fmt.Println(p)
		fmt.Println("================")
	}
	iter.Release()
}
