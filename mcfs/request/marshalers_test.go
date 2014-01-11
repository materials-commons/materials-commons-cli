package request

import (
	"fmt"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

var _ = fmt.Println

func TestRequestMarshaler(t *testing.T) {
	m := NewRequestResponseMarshaler()
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
