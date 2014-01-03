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
	user string
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
			r.user = t.User
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
	case transfer.CreateFile:
		return r.createFile(req)
	case transfer.CreateDir:
		return r.createDir(req)
	case transfer.CreateProject:
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
	return nil
}

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
		return r.badRequest(fmt.Errorf("Bad request data for type %d", req.Type))
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

func (r *ReqHandler) createFile(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.CreateFileReq:
		var _ = t
		return nil
	default:
		return r.badRequest(fmt.Errorf("Bad request data for type %d", req.Type))
	}
}

func (r *ReqHandler) createDir(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.CreateDirReq:
		if r.db.verifyProject(t.ProjectID) {
			return r.createDataDir(t)
		}
		return r.badRequest(fmt.Errorf("Invalid project: %s", t.ProjectID))
	default:
		return r.badRequest(fmt.Errorf("Bad request data for type %d", req.Type))
	}
}

func (db db) verifyProject(projectID string) bool {
	return false
}

func (rh *ReqHandler) createDataDir(req transfer.CreateDirReq) ReqStateFN {
	proj, err := model.GetProject(req.ProjectID, rh.db.session)
	switch {
	case err != nil:
		err = fmt.Errorf("Bad projectID %s", req.ProjectID)
	case proj.Owner != rh.user:
		err = fmt.Errorf("Access to project not allowed")
	default:
		var parent string
		if parent, err = rh.db.getParent(req.Path); err == nil {
			datadir := model.NewDataDir(req.Path, "private", rh.user, parent)
			_, err = r.Table("datadirs").Insert(datadir).RunWrite(rh.db.session)
		}
	}

	if err != nil {
		rh.respError(err)
	} else {
		rh.respContinue()
	}
	return rh.nextCommand()
}

func (db db) getParent(ddirPath string) (string, error) {
	query := r.Table("datadirs").GetAllByIndex("name", ddirPath)
	var d model.DataDir
	err := model.GetRow(query, db.session, &d)
	return d.Id, err
}

func createDataFileDir(dataFileID string) error {
	pieces := strings.Split(dataFileID, "-")
	dirpath := filepath.Join("/mcfs/data/materialscommons", pieces[1][0:2], pieces[1][2:4])
	return os.MkdirAll(dirpath, 0600)
}

func datafileDir(dataFileID string) string {
	pieces := strings.Split(dataFileID, "-")
	return filepath.Join("/mcfs/data/materialscommons", pieces[1][0:2], pieces[1][2:4])
}

func datafilePath(dataFileID string) string {
	return filepath.Join(datafileDir(dataFileID), dataFileID)
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

func (r *ReqHandler) respError(err error) {
	resp := &transfer.Response{
		Type:   transfer.RContinue,
		Status: err,
	}
	r.Encode(resp)
}
