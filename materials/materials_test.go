package main

import (
	"github.com/materials-commons/materials"
	"testing"
)

func TestConvertProjects(t *testing.T) {
	u, _ := materials.NewUserFrom("../test_data/conversion")
	materials.ConfigInitialize(u)
	convertProjects()

	// Make sure the conversion went correctly
	projectDB, err := materials.OpenProjectDB("../test_data/conversion/.materials/projectdb")

	if err != nil {
		t.Fatalf("Unable to open projectdb %s", err.Error())
	}

	for _, project := range projectDB.Projects() {
		switch {
		case project.Name == "proj1a":
			verify(project, "/tmp/proj1a", "Unloaded", t)
		case project.Name == "proj 2a":
			verify(project, "/tmp/proj 2a", "Loaded", t)
		default:
			t.Fatalf("Unexpected project %#v\n", project)
		}
	}
}

func verify(project materials.Project, path, status string, t *testing.T) {
	if project.Path != path {
		t.Fatalf("Paths don't match, expected %s, got %s\n", project.Path, path)
	}

	if project.Status != status {
		t.Fatalf("Status don't match, expected %s, got %s\n", project.Status, status)
	}
}
