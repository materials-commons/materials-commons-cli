package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

var _ = fmt.Println

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
	model.Delete("datadirs", createdId, session)

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
