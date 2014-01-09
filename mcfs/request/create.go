package request

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"path/filepath"
	"strings"
)

func (h *ReqHandler) createProject(req transfer.CreateProjectReq) ReqStateFN {
	switch {
	case !validProjectName(req.Name):
		return h.badRequestNext(fmt.Errorf("Invalid project name %s", req.Name))
	case h.db.projectExists(req.Name, h.user):
		return h.badRequestNext(fmt.Errorf("Project %s exists", req.Name))
	default:
		projectId, datadirId, err := h.db.createProject(req.Name, h.user)
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

func (h *ReqHandler) createFile(req transfer.CreateFileReq) ReqStateFN {
	if err := h.db.validCreateFileReq(req, h.user); err != nil {
		return h.badRequestNext(err)
	}

	df := model.NewDataFile(req.Name, "private", h.user)
	df.DataDirs = append(df.DataDirs, req.DataDirID)
	rv, err := r.Table("datafiles").Insert(df).RunWrite(h.db.session)
	if err != nil {
		return h.badRequestNext(err)
	}

	if rv.Inserted == 0 {
		return h.badRequestNext(fmt.Errorf("Unable to insert datafile"))
	}
	datafileId := rv.GeneratedKeys[0]

	// TODO: Eliminate an extra query to look up the DataDir
	// when we just did during verification.
	datadir, _ := model.GetDataDir(req.DataDirID, h.db.session)
	datadir.DataFiles = append(datadir.DataFiles, datafileId)

	// TODO: Really should check for errors here. What do
	// we do? The database could get out of sync. Maybe
	// need a way to update partially completed items by
	// putting into a log? Ugh...
	r.Table("datadirs").Update(datadir).RunWrite(h.db.session)
	createResp := transfer.CreateResp{
		ID: datafileId,
	}
	h.respOk(createResp)
	return h.nextCommand()
}

func (db db) validCreateFileReq(fileReq transfer.CreateFileReq, user string) error {
	proj, err := model.GetProject(fileReq.ProjectID, db.session)
	if err != nil {
		return fmt.Errorf("Unknown project id %s", fileReq.ProjectID)
	}

	if proj.Owner != user {
		return fmt.Errorf("User %s is not owner of project %s", user, proj.Name)
	}

	datadir, err := model.GetDataDir(fileReq.DataDirID, db.session)
	if err != nil {
		return fmt.Errorf("Unknown datadir Id %s", fileReq.DataDirID)
	}

	if !db.datadirInProject(datadir.Id, proj.Id) {
		return fmt.Errorf("Datadir %s not in project %s", datadir.Name, proj.Name)
	}

	if db.datafileExistsInDataDir(fileReq.DataDirID, fileReq.Name) {
		return fmt.Errorf("Datafile %s already exists in datadir %s", fileReq.Name, datadir.Name)
	}

	return nil
}

type Project2Datadir struct {
	Id        string `gorethink:"id,omitempty"`
	ProjectID string `gorethink:"project_id"`
	DataDirID string `gorethink:"datadir_id"`
}

func (db db) datadirInProject(datadirId, projectId string) bool {
	query := r.Table("project2datadir").GetAllByIndex("datadir_id", datadirId)
	var p2d Project2Datadir
	err := model.GetRow(query, db.session, &p2d)
	switch {
	case err != nil:
		return false
	case p2d.ProjectID != projectId:
		return false
	default:
		return true
	}
}

func (db db) datafileExistsInDataDir(datadirID, datafileName string) bool {
	rows, err := r.Table("datafiles").GetAllByIndex("name", datafileName).Run(db.session)
	if err != nil {
		return true // don't know if it exists or not
	}
	defer rows.Close()

	for rows.Next() {
		var df model.DataFile
		rows.Scan(&df)
		for _, ddirID := range df.DataDirs {
			if datadirID == ddirID {
				return true
			}
		}
	}
	return false
}

func (r *ReqHandler) createDir(req transfer.CreateDirReq) ReqStateFN {
	if r.db.verifyProject(req.ProjectID, r.user) {
		return r.createDataDir(req)
	}
	return r.badRequestNext(fmt.Errorf("Invalid project: %s", req.ProjectID))
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
			var wr r.WriteResponse
			wr, err = r.Table("datadirs").Insert(datadir).RunWrite(rh.db.session)
			if err == nil && wr.Inserted > 0 {
				p2d := Project2Datadir{
					ProjectID: req.ProjectID,
					DataDirID: datadir.Id,
				}
				r.Table("project2datadir").Insert(p2d).RunWrite(rh.db.session)
			}
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
