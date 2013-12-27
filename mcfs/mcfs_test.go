package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/gohandy/rethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

var _ = fmt.Println

var (
	session, _ = r.Connect(map[string]interface{}{
		"address":  "localhost:30815",
		"database": "materialscommons",
	})

	h = transfer.StartHeader{
		ProjectID: "abc123",
		User:      "gtarcea@umich.edu",
		ApiKey:    "472abe203cd411e3a280ac162d80f1bf",
	}

	c = transfer.Command{
		Header: h,
	}

	ch = &commandHandler{
		Command: &c,
		db:      rethink.NewDB(session),
	}
)

func TestValidApiKey(t *testing.T) {
	if !ch.validApiKey() {
		t.Fatalf("Apikey invalid, should have been valid: %s\n", ch.Header.ApiKey)
	}

	ch.Header.User = "doesnot-exist@nosuch.com"
	if ch.validApiKey() {
		t.Fatalf("Apikey check passed for invalid user: %s\n", ch.Header.User)
	}

	ch.Header.User = "gtarcea@umich.edu"
	ch.Header.ApiKey = "abc123"
	if ch.validApiKey() {
		t.Fatalf("Apikey check should have failed: %s\n", ch.Header.ApiKey)
	}
}

func TestHasAccess(t *testing.T) {
	user := "gtarcea@umich.edu"
	owner := "mcfada@umich.edu"
	session := ch.db.Session
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
	rv, err := r.Table("usergroups").Insert(ug).RunWrite(ch.db.Session)
	if err != nil {
		t.Fatalf("Unable to create new usergroup")
	}
	id := rv.GeneratedKeys[0]
	defer deleteItem(id, "usergroups", ch.db.Session)

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
