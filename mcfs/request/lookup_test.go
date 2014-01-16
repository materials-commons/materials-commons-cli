package request

import (
	"fmt"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

var _ = fmt.Println

var dataDirTests = []struct {
	field    string
	value    string
	errorNil bool
	comment  string
}{
	{"id", "abc123", false, "No such id"},
	{"id", "gtarcea@umich.edu$Test_Proj", true, "Existing with permissions"},
	{"id", " mcfada@umich.edu$Synthetic Tooth_Pics for update_9-23-13", false, "Existing without permission"},
	{"blah", "blah", false, "No such field"},
}

func TestLookupDataDir(t *testing.T) {
	h := NewReqHandler(nil, session, "")
	h.user = "gtarcea@umich.edu"
	for _, test := range dataDirTests {
		req := &transfer.LookupReq{
			Field: test.field,
			Value: test.value,
			Type:  "datadir",
		}
		v, err := h.lookup(req)
		switch {
		case err != nil && test.errorNil:
			// Expected error to be nil
			t.Fatalf("Expected error to be nil for test %s, err %s", test.comment, err)
		case err == nil && !test.errorNil:
			// Expected error not to be nil
			t.Fatalf("Expected err != nil for test %s", test.comment)
		default:
			fmt.Printf("%#v\n", v)
		}
	}
}

func TestLookupDataFile(t *testing.T) {

}

func TestLookupProject(t *testing.T) {

}

func TestLookupInvalidItem(t *testing.T) {

}
