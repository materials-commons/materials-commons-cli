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
	for reqStateFN := r.startState; reqStateFN != nil; {
		reqStateFN = reqStateFN()
	}
}

func (r *ReqHandler) req() transfer.Request {
	req := transfer.Request{}
	if err := r.Decode(&req); err != nil {
		// Bad Request create request that will jump to the end.
		return transfer.Request{
			Type: transfer.Error,
		}
	}
	return req
}

func (r *ReqHandler) startState() ReqStateFN {
	req := r.req()
	switch req.Type {
	case transfer.Login:
		return r.login(req)
	default:
		r.badRequest(fmt.Errorf("Bad state change %d\n", req.Type))
		return nil
	}
}

func (r *ReqHandler) login(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.LoginReq:
		if r.validLogin(t.User, t.ApiKey) {
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
		Type: transfer.RError,
		Status: err,
	}
	r.Encode(resp)
	return nil
}

func (r *ReqHandler) validLogin(user, apikey string) bool {
	return false
}

func (r *ReqHandler) respContinue() {
	resp := &transfer.Response{
		Type: transfer.RContinue,
		Status: nil,
	}
	r.Encode(resp)
}

func (r *ReqHandler) nextCommand() ReqStateFN {
	req := transfer.Request{}
	r.Decode(&req)
	switch req.Type {
	case transfer.Upload:
	case transfer.Download:
	case transfer.Move:
	case transfer.Delete:
	case transfer.Stat:
	default:
		return r.badRequest(fmt.Errorf("Bad request in NextCommand: %d", req.Type))
	}
}
