package request

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"path/filepath"
	"strings"
)

func (r *ReqHandler) createFile(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.CreateFileReq:
		var _ = t
		return nil
	default:
		return r.badRequestNext(fmt.Errorf("3 Bad request data for type %d", req.Type))
	}
}

func (r *ReqHandler) createDir(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.CreateDirReq:
		if r.db.verifyProject(t.ProjectID, r.user) {
			return r.createDataDir(t)
		}
		return r.badRequestNext(fmt.Errorf("Invalid project: %s", t.ProjectID))
	default:
		return r.badRequestNext(fmt.Errorf("4 Bad request data for type %d", req.Type))
	}
}

func (db db) verifyProject(projectID, user string) bool {
	project, err := model.GetProject(projectID, db.session)
	switch {
	case err != nil:
		return false
	case project.Owner != user:
		return false
	default:
		return true
	}
}

func (rh *ReqHandler) createDataDir(req transfer.CreateDirReq) ReqStateFN {
	var datadir model.DataDir
	proj, err := model.GetProject(req.ProjectID, rh.db.session)
	switch {
	case err != nil:
		err = fmt.Errorf("Bad projectID %s", req.ProjectID)
	case proj.Owner != rh.user:
		err = fmt.Errorf("Access to project not allowed")
	case !rh.db.validDirPath(proj.Name, req.Path):
		err = fmt.Errorf("Invalid directory path %s", req.Path)
	default:
		var parent string
		if parent, err = rh.db.getParent(req.Path); err == nil {
			datadir = model.NewDataDir(req.Path, "private", rh.user, parent)
			_, err = r.Table("datadirs").Insert(datadir).RunWrite(rh.db.session)
		}
	}

	if err != nil {
		rh.respError(err)
	} else {
		resp := transfer.CreateResp{
			ID: datadir.Id,
		}
		rh.respOk(resp)
	}
	return rh.nextCommand()
}

func (db db) validDirPath(projName, dirPath string) bool {
	slash := strings.Index(dirPath, "/")
	switch {
	case slash == -1:
		return false
	case projName != dirPath[:slash]:
		return false
	default:
		return true
	}
}

func (db db) getParent(ddirPath string) (string, error) {
	parent := filepath.Dir(ddirPath)
	query := r.Table("datadirs").GetAllByIndex("name", parent)
	var d model.DataDir
	err := model.GetRow(query, db.session, &d)
	if err != nil {
		return "", fmt.Errorf("No parent for %s", ddirPath)
	}
	return d.Id, nil
}
