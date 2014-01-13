package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"os"
	"testing"
)

var _ = fmt.Println

func TestUploadCasesFile(t *testing.T) {
	// Test New File
	h := NewReqHandler(nil, session, "")
	h.user = "gtarcea@umich.edu"

	// Test bad upload with non existant DataFileID
	uploadReq := transfer.UploadReq{
		DataFileID: "does not exist",
		Size:       6,
		Checksum:   "abc123",
	}

	resp, err := h.upload(&uploadReq)
	if err == nil {
		t.Fatalf("Upload req succeeded with a non existant datafile id")
	}

	// Test create file and then upload
	createFileRequest := transfer.CreateFileReq{
		ProjectID: "c33edab7-a65f-478e-9fa6-9013271c73ea",
		DataDirID: "gtarcea@umich.edu$Test_Proj_6111_Aluminum_Alloys_Data",
		Name:      "testfile1.txt",
		Size:      6,
		Checksum:  "abc123",
	}

	createResp, _ := h.createFile(&createFileRequest)
	createdId := createResp.ID

	uploadReq.DataFileID = createdId

	resp, err = h.upload(&uploadReq)
	if err != nil {
		t.Fatalf("Failed to start upload on a valid file %s", err)
	}

	if resp.DataFileID != createdId {
		t.Fatalf("Upload created a new version expected id %s, got %s", createdId, resp.DataFileID)
	}

	if resp.Offset != 0 {
		t.Fatalf("Upload asking for offset different than 0 (%d)", resp.Offset)
	}

	// Test create and then upload with size larger
	uploadReq.Size = 7
	resp, err = h.upload(&uploadReq)
	if err == nil {
		t.Fatalf("Upload with different size should have failed")
	}

	// Test create and then upload with size smaller
	uploadReq.Size = 5
	resp, err = h.upload(&uploadReq)
	if err == nil {
		t.Fatalf("Upload with different size should have failed")
	}

	// Test create and then upload with different checksum
	uploadReq.Size = 6
	uploadReq.Checksum = "def456"
	resp, err = h.upload(&uploadReq)
	if err == nil {
		t.Fatalf("Upload with different checksum should have failed")
	}

	// Test create and then upload with different size and checksum
	uploadReq.Size = 7
	uploadReq.Checksum = "def456"
	resp, err = h.upload(&uploadReq)
	if err == nil {
		t.Fatalf("Upload with different checksum should have failed")
	}

	// Test Existing without permissions
	h.user = "mcfada@umich.edu"
	uploadReq.Size = 6
	uploadReq.Checksum = "abc123"
	resp, err = h.upload(&uploadReq)
	if err == nil {
		t.Fatalf("Allowing upload when user doesn't have permission")
	}

	// Test interrupted transfer
	h.mcdir = "/tmp/mcdir"
	h.user = "gtarcea@umich.edu"
	os.MkdirAll(h.mcdir, 0777)
	w, err := datafileOpen(h.mcdir, createdId, 0)
	w.Write([]byte("Hello"))
	w.(*os.File).Sync()

	resp, err = h.upload(&uploadReq)
	if err != nil {
		t.Fatalf("Restart interrupted failed")
	}
	if resp.Offset != 5 {
		t.Fatalf("Offset computation incorrect")
	}
	if resp.DataFileID != createdId {
		t.Fatalf("Tried to create a new datafile id for an interrupted transfer")
	}

	// Test new version with previous interrupted
	uploadReq.Size = 8
	uploadReq.Checksum = "def456"
	resp, err = h.upload(&uploadReq)
	if err == nil {
		t.Fatalf("Allowed to create a new version when a previous version hasn't completed upload")
	}

	// Test new version when previous version has completed the upload
	w.Write([]byte("s")) // Get file to correct size to complete upload
	w.Close()
	resp, err = h.upload(&uploadReq)
	if err != nil {
		t.Fatalf("Cannot create new version of file already uploaded %s", err)
	}

	if resp.DataFileID == createdId {
		t.Fatalf("New ID was not assigned for new version of file")
	}

	if resp.Offset != 0 {
		t.Fatalf("Uploading new version offset should be 0 not %d", resp.Offset)
	}

	fmt.Println("Deleting datafile id", createdId)
	model.Delete("datafiles", createdId, session)
	os.RemoveAll("/tmp/mcdir")
}

func TestUploadNewFile(t *testing.T) {

}
