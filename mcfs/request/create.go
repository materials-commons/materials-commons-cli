package request

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
)

func (r *ReqHandler) createFile(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.CreateFileReq:
		var _ = t
		return nil
	default:
		return r.badRequest(fmt.Errorf("3 Bad request data for type %d", req.Type))
	}
}

func (r *ReqHandler) createDir(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.CreateDirReq:
		if r.db.verifyProject(t.ProjectID) {
			return r.createDataDir(t)
		}
		return r.badRequest(fmt.Errorf("Invalid project: %s", t.ProjectID))
	default:
		return r.badRequest(fmt.Errorf("4 Bad request data for type %d", req.Type))
	}
}

func (db db) verifyProject(projectID string) bool {
	return false
}

func (rh *ReqHandler) createDataDir(req transfer.CreateDirReq) ReqStateFN {
	proj, err := model.GetProject(req.ProjectID, rh.db.session)
	switch {
	case err != nil:
		err = fmt.Errorf("Bad projectID %s", req.ProjectID)
	case proj.Owner != rh.user:
		err = fmt.Errorf("Access to project not allowed")
	default:
		var parent string
		if parent, err = rh.db.getParent(req.Path); err == nil {
			datadir := model.NewDataDir(req.Path, "private", rh.user, parent)
			_, err = r.Table("datadirs").Insert(datadir).RunWrite(rh.db.session)
		}
	}

	if err != nil {
		rh.respError(err)
	} else {
		rh.respContinue()
	}
	return rh.nextCommand()
}

func (db db) getParent(ddirPath string) (string, error) {
	query := r.Table("datadirs").GetAllByIndex("name", ddirPath)
	var d model.DataDir
	err := model.GetRow(query, db.session, &d)
	return d.Id, err
}
