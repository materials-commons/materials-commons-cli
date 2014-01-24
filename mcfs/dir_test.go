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
	{"c33edab7-a65f-478e-9fa6-9013271c73ea", "Test_Proj", "/tmp/abc.txt", false, "Valid project bad path"},
	{"c33edab7-a65f-478e-9fa6-9013271c73ea", "Test_Proj", "Test_Proj/abc.txt", true, "Valid project path starts with project"},
	{"does not exist", "Test_Proj", "Test_Proj/abc.txt", false, "Valid project path bad project id"},
	{"c33edab7-a65f-478e-9fa6-9013271c73ea", "Test_Proj", "/tmp/Test_Proj/abc.txt", true, "Valid project full path containing project name"},
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

	// Test CreateDir with wrong path

	// Test CreateDir with complete path, as opposed to just the Project Portion

	// Test CreateDir with just the Project Portion of the path

	// Test CreateDir with a bad project ID

	// Test CreateDir with a good project ID and wrong path (can't test right now)
}
