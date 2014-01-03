package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"testing"
)

var _ = fmt.Println

var (
	session, _ = r.Connect(map[string]interface{}{
		"address":  "localhost:30815",
		"database": "materialscommons",
	})
)

func TestHasAccess(t *testing.T) {
	user := "gtarcea@umich.edu"
	owner := "mcfada@umich.edu"
	// Test empty table different user
	if ownerGaveAccessTo(owner, "someuser@umich.edu", session) {
		t.Fatalf("Access passed should have failed with empty usergroups table")
	}

	//Test empty table same user
	if !ownerGaveAccessTo("gtarcea@umich.edu", "gtarcea@umich.edu", session) {
		t.Fatalf("Access failed when user is also the user")
	}

	ug := model.NewUserGroup("mcfada@umich.edu", "tgroup1")
	ug.Users = append(ug.Users, "gtarcea@umich.edu")
	rv, err := r.Table("usergroups").Insert(ug).RunWrite(session)
	if err != nil {
		t.Fatalf("Unable to create new usergroup")
	}
	id := rv.GeneratedKeys[0]
	defer deleteItem(id, "usergroups", session)

	// Test user who should have access
	if !ownerGaveAccessTo(owner, user, session) {
		t.Fatalf("gtarcea@umich.edu should have had access")
	}

	// Test user who doesn't have access
	if ownerGaveAccessTo(owner, "nouser@umich.edu", session) {
		t.Fatalf("nouser@umich.edu should not have access")
	}
}

func deleteItem(id, table string, session *r.Session) {
	fmt.Printf("Deleting id %s from table %s\n", id, table)
	r.Table(table).Get(id).Delete().RunWrite(session)
}