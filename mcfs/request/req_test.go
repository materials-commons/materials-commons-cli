package request

import (
	"fmt"
	"testing"
	"io"
	"github.com/materials-commons/materials/transfer"
)

var _ = fmt.Println

func TestReq(t *testing.T) {
	m := NewRequestResponseMarshaler()
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
	if err := m.Marshal(&request); err != nil {
		t.Fatalf("Marshal failed")
	}
	val := h.req()
	switch val.(type) {
	case transfer.LoginReq:
	default:
		t.Fatalf("req returned wrong type %T", val)
	}
}
