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
	session *r.Session
	user    string
	mcdir   string
	marshaling.MarshalUnmarshaler
	badRequestCount int
}

func NewReqHandler(m marshaling.MarshalUnmarshaler, session *r.Session, mcdir string) *ReqHandler {
	return &ReqHandler{
		session:            session,
		MarshalUnmarshaler: m,
		mcdir:              mcdir,
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
	var resp interface{}
	var err error
	request := h.req()
	switch req := request.(type) {
	case transfer.LoginReq:
		resp, err = h.login(&req)
		if err != nil {
			return h.badRequestRestart(err)
		}
		h.respOk(resp)
		return h.nextCommand
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

func (h *ReqHandler) nextCommand() ReqStateFN {
	var err error
	var resp interface{}

	request := h.req()
	switch req := request.(type) {
	case transfer.UploadReq:
		var respUpload *transfer.UploadResp
		respUpload, err = h.upload(&req)
		if err == nil {
			return h.uploadLoop(respUpload)
		}
	case transfer.CreateFileReq:
		resp, err = h.createFile(&req)
	case transfer.CreateDirReq:
		resp, err = h.createDir(&req)
	case transfer.CreateProjectReq:
		resp, err = h.createProject(&req)
	case transfer.DownloadReq:
	case transfer.MoveReq:
	case transfer.DeleteReq:
	case transfer.LookupReq:
		resp, err = h.lookup(&req)
	case transfer.LogoutReq:
		resp, err = h.logout(&req)
		h.sendResp(resp, err)
		return h.startState
	case transfer.StatReq:
		resp, err = h.stat(&req)
	case transfer.CloseReq:
		return nil
	case transfer.IndexReq:
	default:
		h.badRequestCount = h.badRequestCount + 1
		return h.badRequestNext(fmt.Errorf("Bad request %T", req))
	}

	h.sendResp(resp, err)
	return h.nextCommand
}

func (h *ReqHandler) sendResp(resp interface{}, err error) {
	if err != nil {
		h.respError(err)
	} else {
		h.respOk(resp)
	}
}

func (h *ReqHandler) respOk(respData interface{}) {
	resp := &transfer.Response{
		Type: transfer.ROk,
		Resp: respData,
	}
	err := h.Marshal(resp)
	var _ = err
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
