package request

import (
	"encoding/gob"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/transfer"
	"net"
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
	var req transfer.Request
	err := r.Decode(&req)
	fmt.Printf("Got request %#v\n", req)
	switch {
	case err != nil:
		fmt.Println("   Error:", err)
		req.Type = transfer.Close
	case !transfer.ValidType(req.Type):
		req.Type = transfer.Error
	}
	return req
}

func (r *ReqHandler) startState() ReqStateFN {
	req := r.req()
	switch req.Type {
	case transfer.Login:
		return r.login(req)
	case transfer.Close:
		return nil
	default:
		r.respError(fmt.Errorf("Bad state"))
		return r.badRequest(fmt.Errorf("Bad state change %d\n", req.Type))
	}
}

func (r *ReqHandler) badRequest(err error) ReqStateFN {
	fmt.Println(err)
	resp := &transfer.Response{
		Type:   transfer.RError,
		Status: err.Error(),
	}
	r.Encode(resp)
	return r.startState
}

func (r *ReqHandler) respContinue() {
	fmt.Println("respContinue")
	resp := &transfer.Response{
		Type: transfer.RContinue,
	}
	r.Encode(resp)
}

func (r *ReqHandler) nextCommand() ReqStateFN {
	req := transfer.Request{}
	r.Decode(&req)
	fmt.Printf("req = %#v\n", req)
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
	case transfer.Close:
		return nil
	default:
		return r.badRequest(fmt.Errorf("2 Bad request in NextCommand: %d", req.Type))
	}
	return nil
}

func (r *ReqHandler) respError(err error) {
	resp := &transfer.Response{
		Type:   transfer.RContinue,
		Status: err.Error(),
	}
	r.Encode(resp)
}
