package request

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (r *ReqHandler) upload(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.UploadReq:
		offset, err := r.validateUploadReq(t)
		if err != nil {
			return r.badRequest(err)
		}
		r.respUpload(offset, t.DataFileID)
		return r.uploadLoop(t.DataFileID)
	default:
		return r.badRequest(fmt.Errorf("6 Bad request data for type %d", req.Type))
	}
}

func (r *ReqHandler) validateUploadReq(req transfer.UploadReq) (offset int64, err error) {
	offset = -1
	df, err := model.GetDataFile(req.DataFileID, r.db.session)
	switch {
	case err != nil:
		return offset, err
	case req.Checksum != df.Checksum:
		return offset, fmt.Errorf("Checksums don't match")
	default:
		sinfo, err := os.Stat(datafilePath(req.DataFileID))
		if err == nil && sinfo.Size() < df.Size {
			offset = sinfo.Size()
		}
	}
	return offset, err
}

func (r *ReqHandler) respUpload(offset int64, dfid string) {

}

type uploadHandler struct {
	w          io.WriteCloser
	dataFileID string
	nbytes     int64
	*ReqHandler
}

func datafileWrite(w io.WriteCloser, bytes []byte) (int, error) {
	return w.Write(bytes)
}

func datafileClose(w io.WriteCloser, dataFileID string, session *r.Session) error {
	// Update datafile in db
	w.Close()
	return nil
}

func datafileOpen(dfid string) (io.WriteCloser, error) {
	err := createDataFileDir(dfid)
	if err != nil {
		return nil, err
	}
	return os.Create(datafilePath(dfid))
}

/*
The following variables define functions for interacting with the datafile. They also
allow these functions to be replaced during testing when the test doesn't really
need to do anything with the datafile.

TODO: Think about creating a type and interface that defines all operations on a
data file, Then pass this interface in to the request handler. That way we can
always replace it for testing or other purposes.
*/
var dfWrite = datafileWrite
var dfClose = datafileClose
var dfOpen = datafileOpen

func (r *ReqHandler) uploadLoop(dfid string) ReqStateFN {
	f, err := dfOpen(dfid)
	if err != nil {
		return nil // return something else
	}
	h := &uploadHandler{
		w:          f,
		dataFileID: dfid,
		nbytes:     0,
		ReqHandler: r,
	}

	return h.upload()
}

func (h *uploadHandler) upload() ReqStateFN {
	req := h.req()
	switch req.Type {
	case transfer.Send:
		switch t := req.Req.(type) {
		case transfer.SendReq:
			if t.DataFileID != h.dataFileID {
				// bad send - error out?
			}
			n, err := dfWrite(h.w, t.Bytes)
			if err != nil {
				// error writing, do something...
			}
			h.nbytes = h.nbytes + int64(n)
		default:
			// What to do here? Probably assume an error and close
			// the connection.
			dfClose(h.w, h.dataFileID, h.db.session)
			return h.badRequest(fmt.Errorf("Bad Request"))
		}
	case transfer.Error:
	case transfer.Logout:
		dfClose(h.w, h.dataFileID, h.db.session)
		return h.startState
	case transfer.Close:
		dfClose(h.w, h.dataFileID, h.db.session)
		return nil
	case transfer.Done:
		dfClose(h.w, h.dataFileID, h.db.session)
		return h.nextCommand()
	default:
		dfClose(h.w, h.dataFileID, h.db.session)
		return h.badRequest(fmt.Errorf("Unknown Request Type"))
	}
	return h.upload()
}

func createDataFileDir(dataFileID string) error {
	pieces := strings.Split(dataFileID, "-")
	dirpath := filepath.Join("/mcfs/data/materialscommons", pieces[1][0:2], pieces[1][2:4])
	return os.MkdirAll(dirpath, 0600)
}
