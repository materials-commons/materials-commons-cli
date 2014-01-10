package request

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/gohandy/marshaling"
	"github.com/materials-commons/materials/transfer"
	"io"
)

const maxBadRequests = 10

type ReqStateFN func() ReqStateFN

type db struct {
	session *r.Session
}

type ReqHandler struct {
	db
	user string
	marshaling.MarshalUnmarshaler
	badRequestCount int
}

func NewReqHandler(m marshaling.MarshalUnmarshaler, session *r.Session) *ReqHandler {
	return &ReqHandler{
		db:                 db{session: session},
		MarshalUnmarshaler: m,
	}
}

func (h *ReqHandler) Run() {
	for reqStateFN := h.startState; reqStateFN != nil; {
		reqStateFN = reqStateFN()
	}
}

type ErrorReq struct{}

func (h *ReqHandler) req() interface{} {
	var req transfer.Request
	if err := h.Unmarshal(&req); err != nil {
		if err == io.EOF {
			return transfer.CloseReq{}
		}
		return ErrorReq{}
	}
	return req.Req
}

func (h *ReqHandler) startState() ReqStateFN {
	request := h.req()
	switch req := request.(type) {
	case transfer.LoginReq:
		return h.login(&req)
	case transfer.CloseReq:
		return nil
	default:
		return h.badRequestRestart(fmt.Errorf("Bad Request %T", req))
	}
}

func (h *ReqHandler) badRequestRestart(err error) ReqStateFN {
	fmt.Println("badRequestRestart:", err)
	h.badRequestCount = h.badRequestCount + 1
	resp := &transfer.Response{
		Type:   transfer.RError,
		Status: err.Error(),
	}
	h.Marshal(resp)
	if h.badRequestCount > maxBadRequests {
		return nil
	}
	return h.startState
}

func (h *ReqHandler) badRequestNext(err error) ReqStateFN {
	fmt.Println("badRequestNext:", err)
	resp := &transfer.Response{
		Type:   transfer.RError,
		Status: err.Error(),
	}
	h.Marshal(resp)
	if h.badRequestCount > maxBadRequests {
		return nil
	}
	return h.nextCommand
}

func (h *ReqHandler) respOk(respData interface{}) {
	resp := &transfer.Response{
		Type: transfer.ROk,
		Resp: respData,
	}
	h.Marshal(resp)
}

func (h *ReqHandler) nextCommand() ReqStateFN {
	request := h.req()
	switch req := request.(type) {
	case transfer.UploadReq:
		return h.upload(&req)
	case transfer.CreateFileReq:
		return h.createFile(&req)
	case transfer.CreateDirReq:
		return h.createDir(&req)
	case transfer.CreateProjectReq:
		return h.createProject(&req)
	case transfer.DownloadReq:
	case transfer.MoveReq:
	case transfer.DeleteReq:
	case transfer.LogoutReq:
		return h.logout(&req)
	case transfer.StatReq:
		return h.stat(&req)
	case transfer.CloseReq:
		return nil
	case transfer.IndexReq:
	default:
		h.badRequestCount = h.badRequestCount + 1
		return h.badRequestNext(fmt.Errorf("Bad request %T", req))
	}
	return nil
}

func (h *ReqHandler) respError(err error) {
	fmt.Println("respError:", err)
	resp := &transfer.Response{
		Type:   transfer.RError,
		Status: err.Error(),
	}
	h.Marshal(resp)
}

func (h *ReqHandler) respFatal(err error) {
	resp := &transfer.Response{
		Type:   transfer.RFatal,
		Status: err.Error(),
	}
	h.Marshal(resp)
}
