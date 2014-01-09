package request

import (
	"encoding/gob"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/transfer"
	"io"
	"net"
)

const maxBadRequests = 10

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
	badRequestCount int
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

type ErrorReq struct{}

func (r *ReqHandler) req() interface{} {
	var req transfer.Request
	if err := r.Decode(&req); err != nil {
		if err == io.EOF {
			return transfer.CloseReq{}
		}
		return ErrorReq{}
	}
	return req.Req
}

func (r *ReqHandler) startState() ReqStateFN {
	request := r.req()
	switch req := request.(type) {
	case transfer.LoginReq:
		return r.login(&req)
	case transfer.CloseReq:
		return nil
	default:
		return r.badRequestRestart(fmt.Errorf("Bad Request %T", req))
	}
}

func (r *ReqHandler) badRequestRestart(err error) ReqStateFN {
	fmt.Println("badRequestRestart:", err)
	r.badRequestCount = r.badRequestCount + 1
	resp := &transfer.Response{
		Type:   transfer.RError,
		Status: err.Error(),
	}
	r.Encode(resp)
	if r.badRequestCount > maxBadRequests {
		return nil
	}
	return r.startState
}

func (r *ReqHandler) badRequestNext(err error) ReqStateFN {
	fmt.Println("badRequestNext:", err)
	resp := &transfer.Response{
		Type:   transfer.RError,
		Status: err.Error(),
	}
	r.Encode(resp)
	if r.badRequestCount > maxBadRequests {
		return nil
	}
	return r.nextCommand
}

func (r *ReqHandler) respOk(respData interface{}) {
	resp := &transfer.Response{
		Type: transfer.ROk,
		Resp: respData,
	}
	r.Encode(resp)
}

func (r *ReqHandler) nextCommand() ReqStateFN {
	request := r.req()
	switch req := request.(type) {
	case transfer.UploadReq:
		return r.upload(&req)
	case transfer.CreateFileReq:
		return r.createFile(&req)
	case transfer.CreateDirReq:
		return r.createDir(&req)
	case transfer.CreateProjectReq:
		return r.createProject(&req)
	case transfer.DownloadReq:
	case transfer.MoveReq:
	case transfer.DeleteReq:
	case transfer.LogoutReq:
		return r.logout(&req)
	case transfer.StatReq:
		return r.stat(&req)
	case transfer.CloseReq:
		return nil
	case transfer.IndexReq:
	default:
		r.badRequestCount = r.badRequestCount + 1
		return r.badRequestNext(fmt.Errorf("Bad request %T", req))
	}
	return nil
}

func (r *ReqHandler) respError(err error) {
	fmt.Println("respError:", err)
	resp := &transfer.Response{
		Type:   transfer.RError,
		Status: err.Error(),
	}
	r.Encode(resp)
}

func (r *ReqHandler) respFatal(err error) {
	resp := &transfer.Response{
		Type:   transfer.RFatal,
		Status: err.Error(),
	}
	r.Encode(resp)
}
