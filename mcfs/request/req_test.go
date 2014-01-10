package request

import (
	"fmt"
	"testing"
	"io"
	"github.com/materials-commons/materials/transfer"
)

func TestReq(t *testing.T) {
	m := NewRequestMarshaler()
	h := NewReqHandler(m, session)

	m.SetError(io.EOF)
	switch h.req().(type) {
	case transfer.CloseReq:
	default:
		t.Fatalf("Wrong type")
	}

	m.SetError(fmt.Errorf(""))
	switch h.req().(type) {
	case ErrorReq:
	default:
		t.Fatalf("Wrong type")
	}

	m.ClearError()
	loginReq := transfer.LoginReq{}
	request := transfer.Request{ loginReq }
	m.Marshal(request)
	val := h.req()
	fmt.Printf("%#v", val)
}


