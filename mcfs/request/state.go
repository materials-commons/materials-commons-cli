package request

import (
	r "github.com/dancannon/gorethink"
	"net"
	"github.com/materials-commons/materials/transfer"
	"encoding/gob"
	"fmt"
)

type ReqStateFN func() ReqStateFN

type ReqHandler struct {
	conn net.Conn
	session *r.Session
	*gob.Decoder
	*gob.Encoder
}

func NewReqHandler(conn net.Conn, session *r.Session) *ReqHandler {
	return &ReqHandler{
		session: session,
		Decoder: gob.NewDecoder(conn),
		Encoder: gob.NewEncoder(conn),
	}
}

func (r *ReqHandler) Run() {
	for reqStateFN := r.StartState; reqStateFN != nil; {
		reqStateFN = reqStateFN()
	}
}

func (r *ReqHandler) StartState() ReqStateFN {
	req := transfer.Request{}
	r.Decode(&req)
	switch req.Type {
	case transfer.Login:
		return r.Login(req)
	default:
		r.BadRequest(fmt.Errorf("Bad state change %d\n", req.Type))
		return nil
	}
}

func (r *ReqHandler) Login(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.LoginReq:
		if r.validLogin(t.User, t.ApiKey) {
			r.Continue()
			return r.NextCommand()
		} else {
			return r.BadRequest(fmt.Errorf("Bad login %s/%s", t.User, t.ApiKey))
		}
	default:
		return r.BadRequest(fmt.Errorf("Bad request data for type %d", req.Type))
	}
}

func (r *ReqHandler) BadRequest(err error) ReqStateFN {
	resp := &transfer.Response{
		Type: transfer.RError,
		Status: err,
	}
	r.Encode(resp)
	return nil
}

func (r *ReqHandler) validLogin(user, apikey string) bool {
	return false
}

func (r *ReqHandler) Continue() {
	resp := &transfer.Response{
		Type: transfer.RContinue,
		Status: nil,
	}
	r.Encode(resp)
}

func (r *ReqHandler) NextCommand() ReqStateFN {
	req := transfer.Request{}
	r.Decode(&req)
	switch req.Type {
	case transfer.Upload:
	case transfer.Download:
	case transfer.Move:
	case transfer.Delete:
	case transfer.Stat:
	default:
		return r.BadRequest(fmt.Errorf("Bad request in NextCommand: %d", req.Type))
	}
}
