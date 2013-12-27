package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/gohandy/rethink"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

func TestValidApiKey(t *testing.T) {
	session, err := r.Connect(map[string]interface{}{
		"address":  "localhost:30815",
		"database": "materialscommons",
	})

	if err != nil {
		fmt.Println(err)
		t.Fatalf("Could not connect to database")
	}

	h := transfer.StartHeader{
		ProjectID: "abc123",
		User:      "gtarcea@umich.edu",
		ApiKey:    "472abe203cd411e3a280ac162d80f1bf",
	}

	c := transfer.Command{
		Header: h,
	}

	ch := &commandHandler{
		Command: &c,
		db:      rethink.NewDB(session),
	}

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
