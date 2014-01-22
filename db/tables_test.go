package db

import (
	"fmt"
	"os"
	"testing"
)

var _ = fmt.Println

func TestCreate(t *testing.T) {
	err := Create("/tmp/sqltest.db")
	defer cleanup("/tmp/sqltest.db")

	if err != nil {
		t.Fatalf("Unable to create database: %s", err)
	}

}

func cleanup(path string) {
	if true {
		return
	}
	os.RemoveAll(path)
}
