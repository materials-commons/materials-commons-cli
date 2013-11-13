package materials

import (
	"fmt"
	"testing"
)

func TestNonExistantUser(t *testing.T) {
	username := "no-such-user-xxx"
	_, err := ProjectsForUser(username)
	if err == nil {
		t.Errorf("Should not have found user '%s'\n", username)
	}
}

func TestProjectsFrom(t *testing.T) {
	projects, err := ProjectsFrom("test_data")
	if err != nil {
		t.Errorf("TsetProjectsFrom failed loading the test_data projects")
	}

	if len(projects.Projects()) != 2 {
		t.Errorf("Number of projects incorrect, it should have been 2: %d", len(projects.Projects()))
	}
}

func TestProjectsFromWithBadDirectory(t *testing.T) {
	projects, err := ProjectsFrom("no-such-directory")
	if err == nil {
		t.Errorf("ProjectFrom should have returned an error")
	}

	if projects != nil {
		t.Errorf("projects should have been nil")
	}
}

func TestReadCorruptedProjectsFile(t *testing.T) {
	projects, err := ProjectsFrom("test_data/corrupted")
	length := len(projects.Projects())

	if err != nil {
		logFail(t, "err should have been nil")
	}

	if length != 2 {
		t.Errorf("Expected corrupted projects to be 2, got %d", length)
	}
}

func TestProjectAddDuplicate(t *testing.T) {
	p, _ := ProjectsFrom("test_data")
	err := p.Add(Project{Name: "proj1", Path: "/tmp"})
	if err == nil {
		t.Fatalf("Duplicate project was added")
	}

	p2, _ := ProjectsFrom("test_data")
	l := len(p2.Projects())
	if l != 2 {
		for _, p := range p2.Projects() {
			fmt.Println(p)
		}
		t.Fatalf("Expected 2 projects, got %d\n", l)
	}
}

func TestProjectAdd(t *testing.T) {
	p, _ := ProjectsFrom("test_data")
	err := p.Add(Project{Name: "new proj", Path: "/tmp"})
	if err != nil {
		logFail(t, "Add failed to add new project")
	}

	l := len(p.Projects())
	if l != 3 {
		logFail(t, "Expected number of projects to be 3, got %d", l)
	}

	p2, _ := ProjectsFrom("test_data")
	l = len(p2.Projects())
	if l != 3 {
		logFail(t, "Expected number of projects to be 3, got %d", l)
	}
}

func TestProjectRemove(t *testing.T) {
	p, _ := ProjectsFrom("test_data")
	err := p.Remove("new proj")
	if err != nil {
		logFail(t, "Remove failed to add new project")
	}

	p2, _ := ProjectsFrom("test_data")
	l := len(p2.Projects())
	if l != 2 {
		logFail(t, "Expected number of projects to be 2, got %d", l)
	}
}

func TestProjectExists(t *testing.T) {
	p, _ := ProjectsFrom("test_data")
	if p.Exists("does-not-exist") {
		t.Fatalf("Found project that doesn't exist\n")
	}

	if !p.Exists("proj1") {
		t.Fatalf("Failed to find project that should exist: proj1")
	}
}

func logFail(t *testing.T, formatString string, args ...interface{}) {
	formatString += "\n"
	t.Errorf(formatString, args)
	t.FailNow()
}
