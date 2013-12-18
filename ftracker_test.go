package materials

import (
	"encoding/json"
	"fmt"
	"os/user"
	"path/filepath"
	"testing"
)

func TestWalkProject(t *testing.T) {
	if true {
		return
	}
	u, _ := user.Current()
	path := filepath.Join(u.HomeDir, "Dropbox/transfers/materialscommons/WE43 Heat Treatments")
	p := Project{
		Name:   "WE43 Heat Treatments",
		Path:   path,
		Status: "Unknown",
	}

	p.Walk()
}

func TestCreatedDb(t *testing.T) {
	if true {
		return
	}

	db, _ := openFileDB("/tmp/project.db")
	defer db.Close()

	iter := db.NewIterator(nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		var p projectFileInfo
		json.Unmarshal(value, &p)
		fmt.Println(key)
		fmt.Println(p)
		fmt.Println("================")
	}
}
