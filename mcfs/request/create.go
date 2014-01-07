package request

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"path/filepath"
	"strings"
)

func (h *ReqHandler) createProject(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.CreateProjectReq:
		switch {
		case !validProjectName(t.Name):
			return h.badRequestNext(fmt.Errorf("Invalid project name %s", t.Name))
		case h.db.projectExists(t.Name, h.user):
			return h.badRequestNext(fmt.Errorf("Project %s exists", t.Name))
		default:
			projectId, datadirId, err := h.db.createProject(t.Name, h.user)
			if err != nil {
				h.respError(err)
			} else {
				resp := transfer.CreateProjectResp{
					ProjectID: projectId,
					DataDirID: datadirId,
				}
				h.respOk(resp)
			}
			return h.nextCommand()
		}
	default:
		return h.badRequestNext(fmt.Errorf("Bad request data for type %s", req.Type))
	}
}

func validProjectName(projectName string) bool {
	i := strings.Index(projectName, "/")
	return i == -1
}

func (db db) projectExists(projectName, user string) bool {
	results, err := r.Table("projects").Filter(r.Row.Field("owner").Eq(user)).
		Filter(r.Row.Field("name").Eq(projectName)).
		Run(db.session)
	if err != nil {
		return true // Error, we don't know if it exists
	}
	defer results.Close()

	return results.Next()
}

func (db db) createProject(projectName, user string) (projectId, datadirId string, err error) {
	datadir := model.NewDataDir(projectName, "private", user, "")
	rv, err := r.Table("datadirs").Insert(datadir).RunWrite(db.session)
	if err != nil {
		return "", "", err
	}
	datadirId = datadir.Id
	project := model.NewProject(projectName, datadirId, user)
	rv, err = r.Table("projects").Insert(project).RunWrite(db.session)
	if err != nil {
		return "", "", err
	}
	return rv.GeneratedKeys[0], datadirId, nil
}

func (r *ReqHandler) createFile(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.CreateFileReq:
		var _ = t
		return nil
	default:
		return r.badRequestNext(fmt.Errorf("Bad request data for type %s", req.Type))
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
		return r.badRequestNext(fmt.Errorf("Bad request data for type %s", req.Type))
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
