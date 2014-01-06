package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
)

func (r *ReqHandler) login(req transfer.Request) ReqStateFN {
	fmt.Println("login")
	switch t := req.Req.(type) {
	case transfer.LoginReq:
		if r.db.validLogin(t.User, t.ApiKey) {
			r.user = t.User
			r.respContinue()
			return r.nextCommand()
		} else {
			return r.badRequest(fmt.Errorf("Bad login %s/%s", t.User, t.ApiKey))
		}
	default:
		return r.badRequest(fmt.Errorf("1 Bad request data for type %d", req.Type))
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

func (r *ReqHandler) logout(req transfer.Request) ReqStateFN {
	r.respContinue()
	return r.startState
}
