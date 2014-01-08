package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	r "github.com/dancannon/gorethink"
	"testing"
)

var _ = fmt.Println
var _ = r.Table

func TestCreateDir(t *testing.T) {
	client := loginTestUser()

	resp := transfer.Response{}
	var request transfer.Request

	// Test valid path

	createDirReq := transfer.CreateDirReq{
		ProjectID: "904886a7-ea57-4de7-8125-6e18c9736fd0",
		Path:      "WE43 Heat Treatments/tdir1",
	}

	request.Type = transfer.CreateDir
	request.Req = createDirReq

	err := client.Encode(&request)
	err = client.Decode(&resp)
	if err != nil {
		t.Fatalf("Decode failed %s", err)
	}

	if resp.Type != transfer.ROk {
		t.Fatalf("Directory create failed with %s", resp.Status)
	}

	createResp := resp.Resp.(transfer.CreateResp)
	createdId := createResp.ID

	// Test existing directory

	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.ROk {
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
	createDirReq.Path = "DIFFERENTPROJECT/tdir1"
	request.Req = createDirReq
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Create dir outside of project succeeded %#v", resp)
	}

	// Test invalid project id
	createDirReq.ProjectID = "abc123"
	createDirReq.Path = "WE43 Heat Treatments/tdir2"
	request.Req = createDirReq
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Create dir with bad project succeeded %#v", resp)
	}

	// Test that fails if subdirs don't exist

	createDirReq = transfer.CreateDirReq{
		ProjectID: "904886a7-ea57-4de7-8125-6e18c9736fd0",
		Path:      "WE43 Heat Treatments/tdir1/tdir2",
	}

	request.Req = createDirReq
	resp = transfer.Response{}

	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Create dir with missing subdirs succeeded %#v", resp)
	}

	// Test sending a CreateDir command with the wrong
	// type of request object.
	request.Req = "Hello world"
	client.Encode(&request)
	resp = transfer.Response{}
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Sent bad req data and didn't get an error")
	}
}

func TestCreateProject(t *testing.T) {
	client := loginTestUser()
	createProjectReq := transfer.CreateProjectReq{
		Name: "TestProject1__",
	}
	resp := transfer.Response{}
	request := transfer.Request{
		Type: transfer.CreateProject,
		Req:  createProjectReq,
	}

	var _ = client
	var _ = resp
	var _ = request

	// Test create new project
	client.Encode(&request)
	client.Decode(&resp)

	createProjectResp := resp.Resp.(transfer.CreateProjectResp)
	projectId := createProjectResp.ProjectID
	datadirId := createProjectResp.DataDirID

	if resp.Type != transfer.ROk {
		t.Fatalf("Unable to create project")
	}
	// Test create existing project
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)

	model.Delete("datadirs", datadirId, session)
	model.Delete("projects", projectId, session)

	if resp.Type != transfer.RError {
		t.Fatalf("Created an existing project - shouldn't be able to")
	}

	// Test create project with invalid name
	createProjectReq = transfer.CreateProjectReq{
		Name: "/InvalidName",
	}
	request.Req = createProjectReq
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Created project with Invalid name")
	}

	// Test create project with bad req data
	resp = transfer.Response{}
	request.Req = "Invalid data"
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Sent request with bad data and didn't get an error")
	}
}

func TestCreateFile(t *testing.T) {
	client := loginTestUser()
	request := transfer.Request{
		Type: transfer.CreateFile,
	}
	resp := transfer.Response{}

	var _ = client
	var _ = resp
	var _ = request

	// Test create a valid file
	createFileReq := transfer.CreateFileReq{
		ProjectID: "c33edab7-a65f-478e-9fa6-9013271c73ea",
		DataDirID: "gtarcea@umich.edu$Test_Proj_6111_Aluminum_Alloys_Data",
		Name:      "testfile1.txt",
	}

	request.Req = createFileReq

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
	validProjectID := createFileReq.ProjectID
	createFileReq.ProjectID = "abc123-doesnotexist"
	request.Req = createFileReq
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed creation of file with bad projectid")
	}

	// Test creating with an invalid datadir id
	createFileReq.ProjectID = validProjectID
	createFileReq.DataDirID = "abc123-doesnotexist"
	request.Req = createFileReq
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed creation of file with bad projectid")
	}

	// Test creating with a datadir not in project
	createFileReq.DataDirID = "mcfada@umich.edu$Synthetic Tooth_Presentation_MCubed"
	request.Req = createFileReq
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed creation of file in a datadir not in project")
	}

	// Test with bad request data
	request.Req = "hello world"
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Should have received error when sending bad req")
	}
}
