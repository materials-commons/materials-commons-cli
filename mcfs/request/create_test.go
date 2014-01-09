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
	if true {
		return
	}
	client := loginTestUser()

	resp := transfer.Response{}

	// Test valid path

	request := transfer.CreateDirReq{
		ProjectID: "904886a7-ea57-4de7-8125-6e18c9736fd0",
		Path:      "WE43 Heat Treatments/tdir1",
	}

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
	request.Path = "DIFFERENTPROJECT/tdir1"
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Create dir outside of project succeeded %#v", resp)
	}

	// Test invalid project id
	request.ProjectID = "abc123"
	request.Path = "WE43 Heat Treatments/tdir2"
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Create dir with bad project succeeded %#v", resp)
	}

	// Test that fails if subdirs don't exist

	request.ProjectID = "904886a7-ea57-4de7-8125-6e18c9736fd0"
	request.Path = "WE43 Heat Treatments/tdir1/tdir2"

	resp = transfer.Response{}

	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Create dir with missing subdirs succeeded %#v", resp)
	}
}

func TestCreateProject(t *testing.T) {
	if true {
		return
	}
	client := loginTestUser()
	request := transfer.CreateProjectReq{
		Name: "TestProject1__",
	}
	resp := transfer.Response{}

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
	request.Name = "/InvalidName"
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Created project with Invalid name")
	}
}

func TestCreateFile(t *testing.T) {
	if true {
		return
	}
	client := loginTestUser()
	resp := transfer.Response{}

	// Test create a valid file
	request := transfer.CreateFileReq{
		ProjectID: "c33edab7-a65f-478e-9fa6-9013271c73ea",
		DataDirID: "gtarcea@umich.edu$Test_Proj_6111_Aluminum_Alloys_Data",
		Name:      "testfile1.txt",
	}

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
	validProjectID := request.ProjectID
	request.ProjectID = "abc123-doesnotexist"
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed creation of file with bad projectid")
	}

	// Test creating with an invalid datadir id
	request.ProjectID = validProjectID
	request.DataDirID = "abc123-doesnotexist"
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed creation of file with bad projectid")
	}

	// Test creating with a datadir not in project
	request.DataDirID = "mcfada@umich.edu$Synthetic Tooth_Presentation_MCubed"
	resp = transfer.Response{}
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Allowed creation of file in a datadir not in project")
	}
}
