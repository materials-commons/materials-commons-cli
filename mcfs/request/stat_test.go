package request

import (
	"fmt"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

var _ = fmt.Println

func TestStat(t *testing.T) {
	// Test existing file
	client := loginTestUser()
	resp := transfer.Response{}

	statRequest := transfer.StatReq{
		DataFileID: "01cc4163-8c6f-4832-8c7b-15e34e4368ae",
	}
	request := transfer.Request{ &statRequest }
	client.Encode(&request)
	client.Decode(&resp)
	fmt.Printf("statresp = %#v\n", resp)

	// Test file we don't have access to

	// Test sending bad data
}
