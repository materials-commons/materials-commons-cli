package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
)

func (r *ReqHandler) login(req transfer.LoginReq) ReqStateFN {
	if r.db.validLogin(req.User, req.ApiKey) {
		r.user = req.User
		r.respOk(nil)
		return r.nextCommand()
	} else {
		return r.badRequestRestart(fmt.Errorf("Bad login %s/%s", req.User, req.ApiKey))
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

func (r *ReqHandler) logout(req transfer.LogoutReq) ReqStateFN {
	r.respOk(nil)
	return r.startState
}
