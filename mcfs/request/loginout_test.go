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

var gtarceaLoginRequest = transfer.Request{
	Type: transfer.Login,
	Req:  gtarceaLoginReq,
}

func loginTestUser() *client {
	client := newClient()
	client.Encode(&gtarceaLoginRequest)
	resp := transfer.Response{}
	client.Decode(&resp)
	return client
}

func TestLoginLogout(t *testing.T) {
	client := newClient()
	loginReq := transfer.LoginReq{
		ProjectID: "abc123",
		User:      "gtarcea@umich.edu",
		ApiKey:    "472abe203cd411e3a280ac162d80f1bf",
	}
	req := transfer.Request{
		Type: transfer.Login,
		Req:  loginReq,
	}

	client.Encode(&req)
	resp := transfer.Response{}
	err := client.Decode(&resp)
	if err != nil {
		t.Fatalf("Unable to decode response")
	}

	if resp.Type != transfer.ROk {
		t.Fatalf("Unexpected return %d expected %d", resp.Type, transfer.ROk)
	}
	req.Type = transfer.Logout
	client.Encode(&req)
	client.Decode(&resp)
	if resp.Type != transfer.ROk {
		t.Fatalf("Unexpected return %d expected %d", resp.Type, transfer.ROk)
	}
	loginReq.ApiKey = "abc12356"
	req.Req = loginReq
	req.Type = transfer.Login
	client.Encode(&req)
	client.Decode(&resp)
	if resp.Type != transfer.RError {
		t.Fatalf("Unexpected return %d expected %d", resp.Type, transfer.RError)
	}

	if resp.Status == "" {
		t.Fatalf("Status should have contained a message")
	}
}
