package model

import (
	r "github.com/dancannon/gorethink"
	"testing"
)

var (
	session, _ = r.Connect(map[string]interface{}{
		"address":  "localhost:30815",
		"database": "materialscommons",
	})
)

func TestGetUser(t *testing.T) {
	_, err := GetUser("nosuch@nosuch.com", session)
	if err == nil {
		t.Fatalf("Found non-existant user nosuch@nosuch.com")
	}

	u, err := GetUser("gtarcea@umich.edu", session)
	if err != nil {
		t.Fatalf("Didn't find existing user gtarcea@umich.edu: %s", err.Error())
	}

	if u.ApiKey != "472abe203cd411e3a280ac162d80f1bf" {
		t.Fatalf("ApiKey does not match, got %s", u.ApiKey)
	}
}
