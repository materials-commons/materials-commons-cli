package request

import (
	"fmt"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

func TestRequestMarshaler(t *testing.T) {
	m := NewRequestMarshaler()
	loginReq := transfer.LoginReq{}
	var _ = loginReq
	request := transfer.Request{ 1 }
	fmt.Printf("request = %#v\n", request)
	m.Marshal(&request)
	var d transfer.Request
	err := m.Unmarshal(&d)
	fmt.Println(d)
	fmt.Println(err)
	fmt.Printf("d.Req = %#v\n", d.Req)
}
