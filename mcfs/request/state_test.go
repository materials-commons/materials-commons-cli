package request

import (
	//	"encoding/gob"
	//	"net"
	r "github.com/dancannon/gorethink"
	"io"
	//	"github.com/materials-commons/materials/transfer"
	"testing"
)

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

}

func TestCreate(t *testing.T) {

}
