package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

var _ = fmt.Println

func TestCreateDir(t *testing.T) {
	client := newClient()
	client.Encode(&gtarceaLoginRequest)
	resp := transfer.Response{}
	client.Decode(&resp)

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
	fmt.Println(resp)
}
