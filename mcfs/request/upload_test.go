package request

import (
	"fmt"
	"github.com/materials-commons/materials/transfer"
	"testing"
)

var _ = fmt.Println

func TestUploadExistingFile(t *testing.T) {
	client := loginTestUser()
	resp := transfer.Response{}

	var _ = client
	var _ = resp
}
