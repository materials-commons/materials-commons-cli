package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
)

func (r *ReqHandler) login(req *transfer.LoginReq) (*transfer.LoginResp, error) {
	if r.db.validLogin(req.User, req.ApiKey) {
		r.user = req.User
		return &transfer.LoginResp{}, nil
	} else {
		return nil, fmt.Errorf("Bad login %s/%s", req.User, req.ApiKey)
	}
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

func (r *ReqHandler) logout(req *transfer.LogoutReq) (*transfer.LogoutResp, error) {
	return nil, nil
}
