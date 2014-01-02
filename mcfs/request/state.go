package request

import (
	"encoding/gob"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type ReqStateFN func() ReqStateFN

type db struct {
	session *r.Session
}

type ReqHandler struct {
	conn net.Conn
	db
	*gob.Decoder
	*gob.Encoder
}

func NewReqHandler(conn net.Conn, session *r.Session) *ReqHandler {
	return &ReqHandler{
		db:      db{session: session},
		Decoder: gob.NewDecoder(conn),
		Encoder: gob.NewEncoder(conn),
	}
}

func (r *ReqHandler) Run() {
	for reqStateFN := r.startState; reqStateFN != nil; {
		reqStateFN = reqStateFN()
	}
}

func (r *ReqHandler) req() transfer.Request {
	req := transfer.Request{}
	if err := r.Decode(&req); err != nil || !transfer.ValidType(req.Type) {
		// Return type that will cause state machine to abort with error
		req.Type = transfer.Error
	}
	return req
}

func (r *ReqHandler) startState() ReqStateFN {
	req := r.req()
	switch req.Type {
	case transfer.Login:
		return r.login(req)
	default:
		return r.badRequest(fmt.Errorf("Bad state change %d\n", req.Type))
	}
}

func (r *ReqHandler) login(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.LoginReq:
		if r.db.validLogin(t.User, t.ApiKey) {
			r.respContinue()
			return r.nextCommand()
		} else {
			return r.badRequest(fmt.Errorf("Bad login %s/%s", t.User, t.ApiKey))
		}
	default:
		return r.badRequest(fmt.Errorf("Bad request data for type %d", req.Type))
	}
}

func (r *ReqHandler) badRequest(err error) ReqStateFN {
	resp := &transfer.Response{
		Type:   transfer.RError,
		Status: err,
	}
	r.Encode(resp)
	return nil
}

func (db db) validLogin(user, apikey string) bool {
	u, err := model.GetUser(user, db.session)
	switch {
	case err != nil:
		return false
	case u.ApiKey != apikey:
		return false
	default:
		return true
	}
}

func (r *ReqHandler) respContinue() {
	resp := &transfer.Response{
		Type:   transfer.RContinue,
		Status: nil,
	}
	r.Encode(resp)
}

func (r *ReqHandler) nextCommand() ReqStateFN {
	req := transfer.Request{}
	r.Decode(&req)
	switch req.Type {
	case transfer.Upload:
		return r.upload(req)
	case transfer.RestartUpload:
	case transfer.Create:
	case transfer.Download:
	case transfer.Move:
	case transfer.Delete:
	case transfer.Logout:
		return r.logout(req)
	case transfer.Stat:
		return r.stat(req)
	default:
		return r.badRequest(fmt.Errorf("Bad request in NextCommand: %d", req.Type))
	}
}

func (r *ReqHandler) upload(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.UploadReq:
		dfid, err := r.validateUploadReq(t)
		if err != nil {
			return r.badRequest(err)
		}
		r.respUpload(dfid)
		return r.uploadLoop(dfid)
	default:
		return r.badRequest(fmt.Errorf("Bad request data for type %d", req.Type))
	}
}

func (r *ReqHandler) validateUploadReq(req transfer.UploadReq) (dfid string, err error) {
	return "", nil
}

func (r *ReqHandler) respUpload(dfid string) {

}

type uploadHandler struct {
	w          io.WriteCloser
	dataFileID string
	nbytes     int64
	*ReqHandler
}

func datafileWrite(w io.Writer, bytes []byte) (int, error) {
	return w.Write(bytes)
}

func datafileClose(w io.WriteCloser, dataFileID string, session *r.Session) error {
	// Update datafile in db
	w.Close()
	return nil
}

func datafileOpen(dfid string) (io.WriteCloser, error) {
	path, err := createDataFilePath(dfid)
	if err != nil {
		return nil, err
	}
	return os.Create(path)
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
	req := transfer.Request{}
	h.Decode(&req)
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
	case transfer.Done:
		dfClose(h.w, h.dataFileID, h.db.session)
		return h.nextCommand()
	default:
		dfClose(h.w, h.dataFileID, h.db.session)
		return h.badRequest(fmt.Errorf("Unknown Request Type"))
	}
	return h.upload()
}

func createDataFilePath(dataFileID string) (string, error) {
	pieces := strings.Split(dataFileID, "-")
	dirpath := filepath.Join("/mcfs/data/materialscommons", pieces[1][0:2], pieces[1][2:4])
	os.MkdirAll(dirpath, 0600)
	return filepath.Join(dirpath, dataFileID), nil
}

func (r *ReqHandler) logout(req transfer.Request) ReqStateFN {
	return nil
}

func (r *ReqHandler) stat(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.StatReq:
		df, err := model.GetDataFile(t.DataFileID, r.db.session)
		if err != nil {
			r.badRequest(fmt.Errorf("Unknown id %s", t.DataFileID))
		}
		r.respStat(df)
		return r.nextCommand()
	default:
		return r.badRequest(fmt.Errorf("Bad request data for type %d", req.Type))
	}
}

func (r *ReqHandler) respStat(df *model.DataFile) {
	statResp := &transfer.StatResp{
		DataFileID: df.Id,
		Checksum:   df.Checksum,
		Size:       df.Size,
		Birthtime:  df.Birthtime,
		MTime:      df.MTime,
	}
	resp := &transfer.Response{
		Type:   transfer.RContinue,
		Status: nil,
		Resp:   statResp,
	}
	r.Encode(resp)
}
