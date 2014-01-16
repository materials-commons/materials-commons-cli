package request

import (
	"fmt"
	"github.com/materials-commons/materials/transfer"
	"github.com/materials-commons/materials/util"
	"testing"
)

var _ = fmt.Println

func TestRequestMarshaler(t *testing.T) {
	m := util.NewRequestResponseMarshaler()
	request := transfer.Request{1}
	m.Marshal(&request)
	var d transfer.Request
	if err := m.Unmarshal(&d); err != nil {
		t.Fatalf("Unmarshal failed with error %s", err)
	}

	if d.Req != 1 {
		t.Fatalf("Inner item not being properly saved")
	}
}

func TestChannelMarshaler(t *testing.T) {
	m := util.NewChannelMarshaler()
	go responder(m)
	loginReq := transfer.LoginReq{
		User:   "gtarcea@umich.edu",
		ApiKey: "abc123",
	}
	req := transfer.Request{
		Req: loginReq,
	}

	if true {
		return
	}
	m.Marshal(&req)
	var resp transfer.Response
	m.Unmarshal(&resp)
	fmt.Printf("resp = %#v\n", resp)
	l := resp.Resp.(*transfer.LogoutResp)
	fmt.Printf("l = %#v\n, *l = %#v\n", l, *l)
}

func responder(m *util.ChannelMarshaler) {
	var request transfer.Request
	m.Unmarshal(&request)
	fmt.Printf("request = %#v\n", request)
	logoutResp := transfer.LogoutResp{}
	resp := transfer.Response{
		Type: transfer.ROk,
		Resp: &logoutResp,
	}
	m.Marshal(&resp)
}
