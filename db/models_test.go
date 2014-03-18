package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/materials-commons/materials/db/schema"
	"os"
	"testing"
	"time"
)

var _ = fmt.Println

var tdb *sqlx.DB

func init() {
	var err error
	dbArgs := fmt.Sprintf("file:%s?cached=shared&mode=rwc", "/tmp/sqltest.db")
	db, err := sql.Open("sqlite3", dbArgs)
	if err != nil {
		panic("Couldn't open test db")
	}
	err = schema.Create(db)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create test db: %s", err))
	}
	db.Close()
	tdb, err = sqlx.Open("sqlite3", dbArgs)
	if err != nil {
		panic("Couldn't reopen db under sqlx")
	}
	Use(tdb)
}

func TestProjects(t *testing.T) {
	proj := schema.Project{
		Name: "testproject",
		Path: "/tmp/testproject",
		MCID: "abc123",
	}

	err := Projects.Insert(proj)
	if err != nil {
		t.Fatalf("Insert Project %#v into projects failed %s", proj, err)
	}

	projects := []schema.Project{}
	err = Projects.Select(&projects, "select * from projects")
	if err != nil {
		t.Fatalf("Select of projects failed: %s", err)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected to get back 1 project and instead got back %d", len(projects))
	}

	proj.ID = 1 // Set the id because first entry will have id 1
	if projects[0] != proj {
		t.Fatalf("Inserted proj different than retrieved version: i/r %#v/%#v", proj, projects[0])
	}

	// Test retrieve a single project
	var p schema.Project
	err = Projects.Get(&p, "select * from projects where path=$1", proj.Path)
	if err != nil {
		t.Errorf("Unable to retrieve a single project: %s", err)
	}

	if p != proj {
		t.Errorf("Inserted object different from retrieved object i/r %#v/%#v", proj, p)
	}

	// Test retrieve non existing
	err = Projects.Get(&p, "select * from projects where path=$1", "/does/not/exist")
	if err == nil {
		t.Errorf("Retrieved non existing project got: %#v", p)
	}
}

func TestProjectEvents(t *testing.T) {
	event := schema.ProjectEvent{
		Path:      "/tmp/testproject/abc.txt",
		Event:     "Delete",
		EventTime: time.Now(),
		ProjectID: 1,
	}

	err := ProjectEvents.Insert(event)

	if err != nil {
		t.Fatalf("Insert ProjectEvent %#v into project_events failed %s", event, err)
	}

	events := []schema.ProjectEvent{}
	err = ProjectEvents.Select(&events, "select * from project_events")
	if err != nil {
		t.Fatalf("Select of project_events failed: %s", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected to get back 1 event and instead got back %d", len(events))
	}

	event.ID = 1 // we know the first id in the database
	// nil out times since they won't be equal
	if !event.EventTime.Equal(events[0].EventTime) {
		t.Fatalf("Inserted event time not equal to retrieved i/r %#v/%#v", event, events[0])
	}

	// set times on both ends since we cannot indirectly compare the times
	// through structure comparison
	now := time.Now()
	event.EventTime = now
	events[0].EventTime = now
	if event != events[0] {
		t.Fatalf("Inserted event different than retrieved version: i/r %#v/%#v", event, events[0])
	}
}

func TestDataFiles(t *testing.T) {
	now := time.Now()
	datafile := schema.DataFile{
		MCID:       "def456",
		Name:       "abc.txt",
		Path:       "/tmp/testproject/abc.txt",
		DataDirID:  1,
		ProjectID:  1,
		Size:       10,
		Checksum:   "fa23ed7",
		LastUpload: now,
		MTime:      now,
		Version:    1,
		ParentMCID: "",
		Parent:     0,
	}

	err := DataFiles.Insert(datafile)
	if err != nil {
		t.Fatalf("Insert DataFile %#v into datafiles failed %s", datafile, err)
	}

	datafiles := []schema.DataFile{}
	err = DataFiles.Select(&datafiles, "select * from datafiles")
	if err != nil {
		t.Fatalf("Select of datafiles failed: %s", err)
	}

	if len(datafiles) != 1 {
		t.Fatalf("Expected to get back 1 datafile and instead got back %d", len(datafiles))
	}

	datafile.ID = 1 // we know the first id in the database
	if !datafile.LastUpload.Equal(datafiles[0].LastUpload) {
		t.Fatalf("Inserted datafile upload time not equal to retrieved i/r %#v/%#v", datafile, datafiles[0])
	}

	if !datafile.MTime.Equal(datafiles[0].MTime) {
		t.Fatalf("Inserted datafile mtime not equal to retrieved i/r %#v/%#v", datafile, datafiles[0])
	}

	// set times on both ends since we cannot indirectly compare the times
	// through structure comparison
	datafile.LastUpload = now
	datafile.MTime = now
	datafiles[0].LastUpload = now
	datafiles[0].MTime = now
	if datafile != datafiles[0] {
		t.Fatalf("Inserted datafile different than retrieved version: i/r %#v/%#v", datafile, datafiles[0])
	}

	// Test that trigger fired
	var count int
	if err = DataFiles.QueryRow("select count(project_id) from project2datafile;").Scan(&count); err != nil {
		t.Errorf("Select count on project2datafile failed: %s", err)
	}
	if count != 1 {
		t.Errorf("Expected count of 1 for project2datafile, got %d", count)
	}

	count = 0
	if err = DataFiles.QueryRow("select count(datadir_id) from datadir2datafile;").Scan(&count); err != nil {
		t.Errorf("Select count on datadir2datafile failed: %s", err)
	}

	if count != 1 {
		t.Errorf("Expected count of 1 for datadir2datafile, got %d", count)
	}
}

func TestDataDirs(t *testing.T) {
	datadir := schema.DataDir{
		MCID:       "ghi789",
		ProjectID:  1,
		Name:       "testproject",
		Path:       "/tmp/testproject",
		ParentMCID: "",
		Parent:     0,
	}
	var _ = datadir

	err := DataDirs.Insert(datadir)
	if err != nil {
		t.Fatalf("Insert DataDir %#v into datadirs failed %s", datadir, err)
	}

	datadirs := []schema.DataDir{}
	err = DataDirs.Select(&datadirs, "select * from datadirs")

	if err != nil {
		t.Fatalf("Select of datadirs failed: %s", err)
	}

	if len(datadirs) != 1 {
		t.Fatalf("Expected to get back 1 datadir and instead got back %d", len(datadirs))
	}

	datadir.ID = 1 // we know the first id in the database
	if datadir != datadirs[0] {
		t.Fatalf("Inserted datadir different than retrieved version: i/r %#v/%#v", datadirs, datadirs[0])
	}

	// Test that trigger fired
	var count int
	if err = DataDirs.QueryRow("select count(project_id) from project2datadir;").Scan(&count); err != nil {
		t.Errorf("Select count on project2datadir failed: %s", err)
	}
	if count != 1 {
		t.Errorf("Expected count of 1 for project2datadir, got %d", count)
	}

	defer cleanupMT()
}

func cleanupMT() {
	fmt.Println("cleanupMT")
	tdb.Close()
	os.RemoveAll("/tmp/sqltest.db")
}
