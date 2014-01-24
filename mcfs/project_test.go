package mcfs

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/materials-commons/materials/db"
	"github.com/materials-commons/materials/db/schema"
	"os"
	"testing"
)

var tdb *sqlx.DB

func init() {
	os.RemoveAll("/tmp/sqltest.db")
	var err error
	dbArgs := fmt.Sprintf("file:%s?cached=shared&mode=rwc", "/tmp/sqltest.db")
	cdb, err := sql.Open("sqlite3", dbArgs)
	if err != nil {
		panic("Couldn't open test db")
	}
	schema.Create(cdb)
	cdb.Close()
	tdb, err = sqlx.Open("sqlite3", dbArgs)
	if err != nil {
		panic("Couldn't reopen db under sqlx")
	}
	db.Use(tdb)
}

func TestProjectExistence(t *testing.T) {
	p := schema.Project{
		Name: "TestProject",
		Path: "/tmp/TestProject",
	}
	db.Projects.Insert(p)
	// Test Existing Project
	proj, err := projectByPath("/tmp/TestProject")
	if err != nil {
		t.Errorf("Failed to retrieve existing project")
	}
	p.Id = proj.Id
	if *proj != p {
		t.Errorf("Retrieve project differs from inserted version i/r %#v/%#v", p, proj)
	}

	// Test Non Existing Project
	proj, err = projectByPath("/tmp/TestProject-does-not-exist")
	if err == nil {
		t.Errorf("Successfully retrieve a non existing project")
	}

	// Test Project with Same name but different path (it should be found)
	proj, err = projectByPath("/does/not/exist/TestProject")
	if err != nil {
		t.Errorf("Failed to retrieve existing project")
	}
	if *proj != p {
		t.Errorf("Retrieve project differs from inserted version i/r %#v/%#v", p, proj)
	}
}

func TestUploadNewProject(t *testing.T) {
	// Test large upload
	if true {
		return
	}
	err := c.UploadNewProject("/home/gtarcea/ST1")
	if err != nil {
		t.Errorf("Failed to upload %s", err)
	}
}
