package request

import (
	"encoding/gob"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"net"
	"os"
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
	file       *os.File
	dataFileID string
	nbytes     int64
	*ReqHandler
}

func (r *ReqHandler) uploadLoop(dfid string) ReqStateFN {
	f, err := openDataFile(dfid)
	if err != nil {
		return nil // return something else
	}
	uh := &uploadHandler{
		file:       f,
		dataFileID: dfid,
		nbytes:     0,
		ReqHandler: r,
	}

	return uh.upload()
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
			n, err := h.file.Write(t.Bytes)
			if err != nil {
				// error writing, do something...
			}
			h.nbytes = h.nbytes + int64(n)
		default:
		}
	case transfer.Error:
	case transfer.Logout:
	case transfer.Done:
		h.file.Close()
		// Update datafile in db
	default:
		// close file, update datafile in db, return badRequest()
	}
	return h.upload()
}

func openDataFile(dfid string) (*os.File, error) {
	return nil, nil
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
