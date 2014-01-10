package request

import (
	"fmt"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

var _ = fmt.Println

func TestStat(t *testing.T) {
	m := NewRequestResponseMarshaler()
	h := NewReqHandler(m, session)
	h.user = "gtarcea@umich.edu"

	// We do this so that the call exits without doing any marshal calls
	// so that we can look at the response.
	request := transfer.Request{ transfer.CloseReq{} }
	m.Marshal(&request)
	
	// Test existing file
	
	//client := loginTestUser()
	resp := transfer.Response{}

	statRequest := transfer.StatReq{
		DataFileID: "1a455b46-a560-472e-acec-c96482fd655a",
	}

	h.stat(&statRequest)
	m.Unmarshal(&resp)

	if resp.Type != transfer.ROk {
		t.Fatalf("Bad stat request")
	}

	sinfo := statResp(resp.Resp)
	if len(sinfo.DataDirs) != 1 {
		t.Fatalf("DataDirs length incorrect, expected 1 got %d", len(sinfo.DataDirs))
	}

	if sinfo.DataDirs[0] != "gtarcea@umich.edu$WE43 Heat Treatments_AT 250C_AT 2 hours_Atom probe" {
		t.Fatalf("Datadirs[0] incorrect = %s", sinfo.DataDirs[0])
	}

	if sinfo.Name != "R38_03085-v01_MassSpectrum.csv" {
		t.Fatalf("Name incorrect = %s", sinfo.Name)
	}

	if sinfo.Checksum != "6a600da8fe52310128ba7f193f6bb345" {
		t.Fatalf("Checksum incorrect = %s", sinfo.Checksum)
	}

	if sinfo.Size != 20637765 {
		t.Fatalf("Size incorrect = %d", sinfo.Size)
	}

	// Test file we don't have access to
	statRequest.DataFileID = "01cc4163-8c6f-4832-8c7b-15e34e4368ae"
	h.stat(&statRequest)
	resp = transfer.Response{}
	m.Unmarshal(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Access to file we shouldn't have access to")
	}

	// Test sending bad DataFileID
	statRequest.DataFileID = "idonotexist"
	h.stat(&statRequest)
	resp = transfer.Response{}
	m.Unmarshal(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Succeeded for data file that doesn't exist")
	}
}

func statResp(req interface{}) transfer.StatResp {
	switch t := req.(type) {
	case *transfer.StatResp:
		return *t
	case transfer.StatResp:
		return t
	default:
		return transfer.StatResp{}
	}
}
