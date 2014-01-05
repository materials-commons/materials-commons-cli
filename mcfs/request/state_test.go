package request

import (
	"encoding/gob"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/transfer"
	"io"
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

	//gob.Register(transfer.LoginReq{})
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

func tdatafileOpen(dfid string) (io.WriteCloser, error) {
	return nil, nil
}

func tdatafileClose(w io.WriteCloser, dataFileID string, session *r.Session) error {
	return nil
}

func tdatafileWrite(w io.WriteCloser, bytes []byte) (int, error) {
	return len(bytes), nil
}

func TestLoginLogout(t *testing.T) {
	client := newClient()
	fmt.Println("got client")

	loginReq := transfer.LoginReq{
		ProjectID: "abc123",
		User:      "gtarcea@umich.edu",
		ApiKey:    "472abe203cd411e3a280ac162d80f1bf",
	}
	req := transfer.Request{
		Type: transfer.Login,
		Req:  loginReq,
	}

	fmt.Println("sending req")
	err := client.Encode(&req)
	fmt.Println(err)
	fmt.Println("req sent")
	resp := transfer.Response{}
	fmt.Println("getting resp")
	err = client.Decode(&resp)
	fmt.Println(err)
	fmt.Printf("%#v", resp)
}

func TestCreate(t *testing.T) {

}
