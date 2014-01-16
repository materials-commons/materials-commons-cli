package request

import (
	"github.com/materials-commons/materials/transfer"
	"github.com/materials-commons/materials/model"
	"fmt"
)

type lookupHandler struct {
	*ReqHandler
	ID string
}

func (h *ReqHandler) lookup(req *transfer.LookupReq) (interface{}, error) {
	l := &lookupHandler{
		ReqHandler: h,
		ID: req.ID,
	}
	switch req.EntryType {
	case "project":
		return l.lookupProject()
	case "datafile":
		return l.lookupDataFile()
	case "datadir":
		return l.lookupDataDir()
	default:
		return nil, fmt.Errorf("Unknown entry type %s", req.EntryType)
	}
}

func (l *lookupHandler) lookupProject() (*model.Project, error) {
	proj, err := model.GetProject(l.ID, l.session)
	switch {
	case err != nil:
		return nil, err
	case !OwnerGaveAccessTo(proj.Owner, l.user, l.session):
		return nil, fmt.Errorf("Permission denied")
	default:
		return proj, nil
	}
}

func (l *lookupHandler) lookupDataFile() (*model.DataFile, error) {
	dataFile, err := model.GetDataFile(l.ID, l.session)
	switch {
	case err != nil:
		return nil, err
	case !OwnerGaveAccessTo(dataFile.Owner, l.user, l.session):
		return nil, fmt.Errorf("Permission denied")
	default:
		return dataFile, nil
	}
}

func (l *lookupHandler) lookupDataDir() (*model.DataDir, error) {
	dataDir, err := model.GetDataDir(l.ID, l.session)
	switch {
	case err != nil:
		return nil, err
	case !OwnerGaveAccessTo(dataDir.Owner, l.user, l.session):
		return nil, fmt.Errorf("Permission denied")
	default:
		return dataDir, nil
	}
}
