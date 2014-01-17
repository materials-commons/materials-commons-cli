package request

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/contrib/model"
	"github.com/materials-commons/materials/transfer"
)

func (h *ReqHandler) login(req *transfer.LoginReq) (*transfer.LoginResp, error) {
	if validLogin(req.User, req.ApiKey, h.session) {
		h.user = req.User
		return &transfer.LoginResp{}, nil
	} else {
		return nil, fmt.Errorf("Bad login %s/%s", req.User, req.ApiKey)
	}
}

func validLogin(user, apikey string, session *r.Session) bool {
	u, err := model.GetUser(user, session)
	switch {
	case err != nil:
		return false
	case u.ApiKey != apikey:
		return false
	default:
		return true
	}
}

func (r *ReqHandler) logout(req *transfer.LogoutReq) (*transfer.LogoutResp, error) {
	return &transfer.LogoutResp{}, nil
}
