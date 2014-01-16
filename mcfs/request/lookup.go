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
		if req.Field != "id" {
			return nil, fmt.Errorf("Projects can only be queried by id")
		}
		rql := l.projectRql(req)
		var proj model.Project
		return l.execute(rql, &proj)

	case "datafile":
		rql := l.dataFileRql(req)
		var datafile model.DataFile
		return l.execute(rql, &datafile)

	case "datadir":
		rql := l.dataDirRql(req)
		var datadir model.DataDir
		return l.execute(rql, &datadir)

	default:
		return nil, fmt.Errorf("Unknown entry type %s", req.Type)
	}
}

func (l *lookupHandler) projectRql(req *transfer.LookupReq) r.RqlTerm {
	return r.Table("projects").Get(req.Value)
}

func (l *lookupHandler) dataFileRql(req *transfer.LookupReq) r.RqlTerm {
	if req.Field == "id" {
		return r.Table("datafiles").Get(req.Value)
	} else {
		return r.Table("datadirs").Filter(r.Row.Field("id").Eq(req.LimitToID)).
			OuterJoin(r.Table("datafiles"),
			func(ddirRow, dfRow r.RqlTerm) r.RqlTerm {
				return ddirRow.Field("datafiles").Contains(dfRow.Field("id"))
			}).Zip().Filter(r.Row.Field(req.Field).Eq(req.Value))
	}
}

func (l *lookupHandler) dataDirRql(req *transfer.LookupReq) r.RqlTerm {
	if req.Field == "id" {
		return r.Table("datadirs").Get(req.Value)
	} else {
		return r.Table("project2datadir").Filter(r.Row.Field("project_id").Eq(req.LimitToID)).
			EqJoin("datadir_id", r.Table("datadirs")).Zip().
			Filter(r.Row.Field(req.Field).Eq(req.Value))
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
