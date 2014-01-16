package request

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
)

type lookupHandler struct {
	session *r.Session
	user    string
}

func (h *ReqHandler) lookup(req *transfer.LookupReq) (interface{}, error) {
	l := &lookupHandler{
		session: h.session,
		user:    h.user,
	}

	switch req.Type {
	case "project":
		rql := r.Table("projects").GetAllByIndex(req.Field, req.Value)
		var proj model.Project
		return l.execute(rql, &proj)
	case "datafile":
		rql := r.Table("datafiles").GetAllByIndex(req.Field, req.Value)
		var datafile model.DataFile
		return l.execute(rql, &datafile)
	case "datadir":
		rql := r.Table("datadirs").GetAllByIndex(req.Field, req.Value)
		var datadir model.DataDir
		return l.execute(rql, &datadir)
	default:
		return nil, fmt.Errorf("Unknown entry type %s", req.Type)
	}
}

func (l *lookupHandler) execute(query r.RqlTerm, v interface{}) (interface{}, error) {
	err := model.GetRow(query, l.session, v)
	switch {
	case err != nil:
		return nil, err
	case !l.hasAccess(v):
		return nil, fmt.Errorf("Permission denied")
	default:
		return v, nil
	}
}

func (l *lookupHandler) hasAccess(v interface{}) bool {
	var owner string
	switch t := v.(type) {
	case *model.Project:
		owner = t.Owner
	case *model.DataDir:
		owner = t.Owner
	case *model.DataFile:
		owner = t.Owner
	}
	return OwnerGaveAccessTo(owner, l.user, l.session)
}
