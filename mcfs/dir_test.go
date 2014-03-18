package mcfs

import (
	"fmt"
	"testing"
)

var _ = fmt.Println

var createDirTests = []struct {
	projectID   string
	projectName string
	path        string
	errorNil    bool
	description string
}{
	{"c33edab7-a65f-478e-9fa6-9013271c73ea", "Test_Proj", "/tmp/abc", false, "Valid project bad path"},
	{"c33edab7-a65f-478e-9fa6-9013271c73ea", "Test_Proj", "Test_Proj/abc", true, "Valid project path starts with project"},
	{"does not exist", "Test_Proj", "Test_Proj/abc.txt", false, "Valid project path bad project id"},
	{"c33edab7-a65f-478e-9fa6-9013271c73ea", "Test_Proj", "/tmp/Test_Proj/abc", true, "Valid project full path containing project name"},
	{"c33edab7-a65f-478e-9fa6-9013271c73ea", "Test_Proj", "/tmp/Test_Proj/abc", true, "Valid project full path containing project name"},
}

func TestCreateDir(t *testing.T) {

	for _, test := range createDirTests {
		_, err := c.CreateDir(test.projectID, test.projectName, test.path)
		switch {
		case err != nil && test.errorNil:
			t.Errorf("Expected error to be nil for test %s, err %s", test.description, err)
		case err == nil && !test.errorNil:
			t.Errorf("Expected err != nil for test %s", test.description)
		}
	}

	// Test creating an existing directory
	projID := "c33edab7-a65f-478e-9fa6-9013271c73ea"
	projName := "Test_Proj"
	dirPath := "/tmp/Test_Proj/abc"
	dataDirID, err := c.CreateDir(projID, projName, dirPath)
	if err != nil {
		t.Errorf("Creating a directory that already exists returned wrong error code: %s", err)
	}
	if dataDirID == "" {
		t.Errorf("Creating an existing directory should have returned the id of the already created directory.")
	}
}
