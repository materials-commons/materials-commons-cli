package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"testing"
	"time"
)

var _ = fmt.Println

var tdb *sqlx.DB

func init() {
	Create("/tmp/sqltest.db")
	var err error
	dbArgs := fmt.Sprintf("file:%s?cached=shared&mode=rwc", "/tmp/sqltest.db")
	tdb, err = sqlx.Open("sqlite3", dbArgs)
	if err != nil {
		panic("Couldn't create test db")
	}
	Use(tdb)
}

func TestProjects(t *testing.T) {
	proj := Project{
		Name: "testproject",
		Path: "/tmp/testproject",
		MCId: "abc123",
	}

	err := Projects.Insert(proj)
	if err != nil {
		t.Fatalf("Insert Project %#v into projects failed %s", proj, err)
	}

	projects := []Project{}
	err = Projects.Select(&projects, "select * from projects")
	if err != nil {
		t.Fatalf("Select of projects failed: %s", err)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected to get back 1 project and instead got back %d", len(projects))
	}

	proj.Id = 1 // Set the id because first entry will have id 1
	if projects[0] != proj {
		t.Fatalf("Inserted proj different than retrieved version: i/r %#v/%#v", proj, projects[0])
	}
}

func TestProjectEvents(t *testing.T) {
	event := ProjectEvent{
		Path:      "/tmp/testproject/abc.txt",
		Event:     "Delete",
		EventTime: time.Now(),
		ProjectId: 1,
	}

	err := ProjectEvents.Insert(event)

	if err != nil {
		t.Fatalf("Insert ProjectEvent %#v into project_events failed %s", event, err)
	}

	events := []ProjectEvent{}
	err = ProjectEvents.Select(&events, "select * from project_events")
	if err != nil {
		t.Fatalf("Select of project_events failed: %s", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected to get back 1 event and instead got back %d", len(events))
	}

	event.Id = 1 // we know the first id in the database
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
	datafile := DataFile{
		MCId:       "def456",
		Name:       "abc.txt",
		Path:       "/tmp/testproject/abc.txt",
		DataDirID:  1,
		ProjectID:  1,
		Size:       10,
		Checksum:   "fa23ed7",
		LastUpload: now,
		MTime:      now,
		Version:    1,
		ParentMCId: "",
		Parent:     0,
	}

	err := DataFiles.Insert(datafile)
	if err != nil {
		t.Fatalf("Insert DataFile %#v into datafiles failed %s", datafile, err)
	}

	datafiles := []DataFile{}
	err = DataFiles.Select(&datafiles, "select * from datafiles")
	if err != nil {
		t.Fatalf("Select of datafiles failed: %s", err)
	}

	if len(datafiles) != 1 {
		t.Fatalf("Expected to get back 1 datafile and instead got back %d", len(datafiles))
	}

	datafile.Id = 1 // we know the first id in the database
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
	datadir := DataDir{
		MCId:       "ghi789",
		ProjectID:  1,
		Name:       "testproject",
		Path:       "/tmp/testproject",
		ParentMCId: "",
		Parent:     0,
	}
	var _ = datadir

	err := DataDirs.Insert(datadir)
	if err != nil {
		t.Fatalf("Insert DataDir %#v into datadirs failed %s", datadir, err)
	}

	datadirs := []DataDir{}
	err = DataDirs.Select(&datadirs, "select * from datadirs")

	if err != nil {
		t.Fatalf("Select of datadirs failed: %s", err)
	}

	if len(datadirs) != 1 {
		t.Fatalf("Expected to get back 1 datadir and instead got back %d", len(datadirs))
	}

	datadir.Id = 1 // we know the first id in the database
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
