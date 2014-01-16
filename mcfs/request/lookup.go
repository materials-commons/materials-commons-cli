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
	var rql r.RqlTerm

	switch req.Type {
	case "project":
		if req.Field == "id" {
			rql = r.Table("projects").Get(req.Value)
		} else {
			rql = r.Table("projects").GetAllByIndex(req.Field, req.Value)
		}
		var proj model.Project
		return l.execute(rql, &proj)
	case "datafile":
		if req.Field == "id" {
			rql = r.Table("datafiles").Get(req.Value)
		} else {
			/*
			    selection = list(r.table('datadirs').filter({'owner': user}).outer_join(\
			 44             r.table('datafiles'), lambda ddirrow, drow: ddirrow['datafiles'].contains(drow['id']))\
			 45                      .run(g.conn, time_format='raw'))
			*/
			rql = r.Table("datadirs").Filter(r.Row.Field("id").Eq(req.LimitToID)).
				OuterJoin(r.Table("datafiles"),
				func(ddirRow, dfRow r.RqlTerm) r.RqlTerm {
					return ddirRow.Field("datafiles").Contains(dfRow.Field("id"))
				}).Zip().Filter(r.Row.Field(req.Field).Eq(req.Value))
		}
		//rql := r.Table("datafiles").GetAllByIndex(req.Field, req.Value)
		var datafile model.DataFile
		return l.execute(rql, &datafile)
	case "datadir":
		if req.Field == "id" {
			rql = r.Table("datadirs").Get(req.Value)
		} else {
			rql = r.Table("project2datadir").Filter(r.Row.Field("project_id").Eq(req.LimitToID)).
				EqJoin("datadir_id", r.Table("datadirs")).Zip().
				Filter(r.Row.Field(req.Field).Eq(req.Value))
		}
		//rql := r.Table("datadirs").GetAllByIndex(req.Field, req.Value)
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
