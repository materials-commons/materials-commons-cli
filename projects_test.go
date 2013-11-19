package materials

import (
	"fmt"
	"testing"
)

func TestNonExistantUser(t *testing.T) {
	username := "no-such-user-xxx"
	_, err := ProjectsForUser(username)
	if err == nil {
		t.Fatalf("Should not have found user '%s'\n", username)
	}
}

func TestProjectsFrom(t *testing.T) {
	projects, err := ProjectsFrom("test_data")
	if err != nil {
		t.Fatalf("TestProjectsFrom failed loading the test_data projects\n")
	}

	if len(projects.Projects()) != 2 {
		t.Fatalf("Number of projects incorrect, it should have been 2: %d\n", len(projects.Projects()))
	}
}

func TestProjectsFromWithBadDirectory(t *testing.T) {
	projects, err := ProjectsFrom("no-such-directory")
	if err == nil {
		t.Fatalf("ProjectFrom should have returned an error\n")
	}

	if projects != nil {
		t.Fatalf("projects should have been nil\n")
	}
}

func TestReadCorruptedProjectsFile(t *testing.T) {
	projects, err := ProjectsFrom("test_data/corrupted")
	length := len(projects.Projects())

	if err != nil {
		t.Fatalf("err should have been nil\n")
	}

	if length != 2 {
		t.Fatalf("Expected corrupted projects to be 2, got %d\n", length)
	}
}

func TestProjectAddDuplicate(t *testing.T) {
	p, _ := ProjectsFrom("test_data")
	err := p.Add(Project{Name: "proj1", Path: "/tmp", Status: "Unloaded"})
	if err == nil {
		t.Fatalf("Duplicate project was added\n")
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
	err := p.Add(Project{Name: "new proj", Path: "/tmp", Status: "Unloaded"})
	if err != nil {
		t.Fatalf("Add failed to add new project\n")
	}

	l := len(p.Projects())
	if l != 3 {
		t.Fatalf("Expected number of projects to be 3, got %d\n", l)
	}

	p2, _ := ProjectsFrom("test_data")
	l = len(p2.Projects())
	if l != 3 {
		t.Fatalf("Expected number of projects to be 3, got %d\n", l)
	}
}

func TestProjectRemove(t *testing.T) {
	p, _ := ProjectsFrom("test_data")
	err := p.Remove("new proj")
	if err != nil {
		t.Fatalf("Remove failed to add new project\n")
	}

	p2, _ := ProjectsFrom("test_data")
	l := len(p2.Projects())
	if l != 2 {
		t.Fatalf("Expected number of projects to be 2, got %d\n", l)
	}
}

func TestProjectExists(t *testing.T) {
	p, _ := ProjectsFrom("test_data")
	if p.Exists("does-not-exist") {
		t.Fatalf("Found project that doesn't exist\n")
	}

	if !p.Exists("proj1") {
		t.Fatalf("Failed to find project that should exist: proj1\n")
	}
}

func TestProjectFind(t *testing.T) {
	p, _ := ProjectsFrom("test_data")

	_, found := p.Find("does-not-exist")
	if found {
		t.Fatalf("Found project that does not exist")
	}

	_, found = p.Find("proj1")
	if !found {
		t.Fatalf("Did not find project proj1\n")
	}

	p.Add(Project{Name: "newproj", Path: "/tmp/newproj"})
	_, found = p.Find("newproj")
	if !found {
		t.Fatalf("Did not find added project newproj\n")
	}

	p.Remove("newproj")
	_, found = p.Find("newproj")
	if found {
		t.Fatalf("Found project that was just removed: newproj\n")
	}
}

func TestProjectUpdate(t *testing.T) {
	p, _ := ProjectsFrom("test_data")
	proj, _ := p.Find("proj1")
	proj.Status = "Loaded"
	p.Update(proj)
	proj, _ = p.Find("proj1")
	if proj.Status != "Loaded" {
		t.Fatalf("proj1 status is %s, should have been 'Loaded'", proj.Status)
	}

	p2, _ := ProjectsFrom("test_data")
	proj, _ = p2.Find("proj1")
	if proj.Status != "Loaded" {
		t.Fatalf("proj1 status is %s, should have been 'Loaded'", proj.Status)
	}
}
