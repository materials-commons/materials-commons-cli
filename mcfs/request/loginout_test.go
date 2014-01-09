package request

import (
	"encoding/gob"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/transfer"
	"net"
	"os"
	"testing"
)

var session *r.Session

func init() {
	session, _ = r.Connect(map[string]interface{}{
		"address":  "localhost:30815",
		"database": "materialscommons",
	})
}

type client struct {
	*gob.Encoder
	*gob.Decoder
}

func newClient() *client {
	conn, err := net.Dial("tcp", "localhost:35862")
	if err != nil {
		fmt.Printf("Couldn't connect %s\n", err.Error())
		os.Exit(1)
	}
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	return &client{
		Encoder: encoder,
		Decoder: decoder,
	}
}

var gtarceaLoginReq = transfer.LoginReq{
	User:   "gtarcea@umich.edu",
	ApiKey: "472abe203cd411e3a280ac162d80f1bf",
}

func loginTestUser() *client {
	client := newClient()
	request := transfer.Request{&gtarceaLoginReq}
	client.Encode(&request)
	resp := transfer.Response{}
	client.Decode(&resp)
	return client
}

func TestLoginLogout(t *testing.T) {
	client := newClient()
	loginRequest := transfer.LoginReq{
		ProjectID: "abc123",
		User:      "gtarcea@umich.edu",
		ApiKey:    "472abe203cd411e3a280ac162d80f1bf",
	}

	request := transfer.Request{&loginRequest}

	client.Encode(&request)
	resp := transfer.Response{}
	err := client.Decode(&resp)
	if err != nil {
		t.Fatalf("Unable to decode response")
	}

	if resp.Type != transfer.ROk {
		t.Fatalf("Unexpected return %d expected %d", resp.Type, transfer.ROk)
	}
	requestLogout := transfer.LogoutReq{}
	request.Req = requestLogout
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.ROk {
		t.Fatalf("Unexpected return %d expected %d", resp.Type, transfer.ROk)
	}
	loginRequest.ApiKey = "abc12356"
	request.Req = loginRequest
	client.Encode(&request)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Unexpected return %d expected %d", resp.Type, transfer.RError)
	}

	if resp.Status == "" {
		t.Fatalf("Status should have contained a message")
	}
}
