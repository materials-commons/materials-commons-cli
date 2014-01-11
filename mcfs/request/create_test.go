package request

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

var _ = fmt.Println
var _ = r.Table

func TestCreateDir(t *testing.T) {
	h := NewReqHandler(nil, session)
	h.user = "gtarcea@umich.edu"

	// Test valid path

	createDirRequest := transfer.CreateDirReq{
		ProjectID: "904886a7-ea57-4de7-8125-6e18c9736fd0",
		Path:      "WE43 Heat Treatments/tdir1",
	}

	resp, err := h.createDir(&createDirRequest)

	if err != nil {
		t.Fatalf("Directory create failed with %s", err)
	}

	createdId := resp.ID
	var _ = createdId

	// Test existing directory

	resp, err = h.createDir(&createDirRequest)
	if err != nil {
		t.Fatalf("Create existing directory failed with %#v", resp)
	}

	// Cleanup the created directory
	fmt.Println("Deleting datadir id:", createdId)
	model.Delete("datadirs", createdId, session)
	// Now cleanup the join table
	rv, _ := r.Table("project2datadir").GetAllByIndex("datadir_id", createdId).Delete().RunWrite(session)
	if rv.Deleted != 1 {
		t.Fatalf("Multiple entries in project2datadir matched. There should only have been one: %#v\n", rv)
	}

	// Test path outside of project
	createDirRequest.Path = "DIFFERENTPROJECT/tdir1"
	resp, err = h.createDir(&createDirRequest)
	if err == nil {
		t.Fatalf("Create dir outside of project succeeded %#v", resp)
	}

	// Test invalid project id
	createDirRequest.ProjectID = "abc123"
	createDirRequest.Path = "WE43 Heat Treatments/tdir2"
	resp, err = h.createDir(&createDirRequest)
	if err == nil {
		t.Fatalf("Create dir with bad project succeeded %#v", resp)
	}

	// Test that fails if subdirs don't exist

	createDirRequest.ProjectID = "904886a7-ea57-4de7-8125-6e18c9736fd0"
	createDirRequest.Path = "WE43 Heat Treatments/tdir1/tdir2"

	resp, err = h.createDir(&createDirRequest)
	if err == nil {
		t.Fatalf("Create dir with missing subdirs succeeded %#v", resp)
	}
}

func TestCreateProject(t *testing.T) {
	h := NewReqHandler(nil, session)
	h.user = "gtarcea@umich.edu"

	createProjectRequest := transfer.CreateProjectReq{
		Name: "TestProject1__",
	}

	// Test create new project
	resp, err := h.createProject(&createProjectRequest)

	projectId := resp.ProjectID
	datadirId := resp.DataDirID

	if err != nil {
		t.Fatalf("Unable to create project")
	}

	// Test create existing project
	resp, err = h.createProject(&createProjectRequest)

	// Delete before test so we can cleanup if there is a failure
	model.Delete("datadirs", datadirId, session)
	model.Delete("projects", projectId, session)

	if err == nil {
		t.Fatalf("Created an existing project - shouldn't be able to")
	}

	// Test create project with invalid name
	createProjectRequest.Name = "/InvalidName"
	resp, err = h.createProject(&createProjectRequest)
	if err == nil {
		t.Fatalf("Created project with Invalid name")
	}
}

func TestCreateFile(t *testing.T) {
	client := loginTestUser()
	resp := transfer.Response{}

	// Test create a valid file
	createFileRequest := transfer.CreateFileReq{
		ProjectID: "c33edab7-a65f-478e-9fa6-9013271c73ea",
		DataDirID: "gtarcea@umich.edu$Test_Proj_6111_Aluminum_Alloys_Data",
		Name:      "testfile1.txt",
	}

	request := transfer.Request{&createFileRequest}

	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.ROk {
		t.Fatalf("Creating file failed")
	}
	createResp := resp.Resp.(transfer.CreateResp)
	createdId := createResp.ID

	// Test creating an existing file
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed create of an existing file")
	}

	// Delete created file
	model.Delete("datafiles", createdId, session)

	// Test creating with an invalid project id
	validProjectID := createFileRequest.ProjectID
	createFileRequest.ProjectID = "abc123-doesnotexist"
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed creation of file with bad projectid")
	}

	// Test creating with an invalid datadir id
	createFileRequest.ProjectID = validProjectID
	createFileRequest.DataDirID = "abc123-doesnotexist"
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed creation of file with bad projectid")
	}

	// Test creating with a datadir not in project
	createFileRequest.DataDirID = "mcfada@umich.edu$Synthetic Tooth_Presentation_MCubed"
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed creation of file in a datadir not in project")
	}
}
