package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
)

func (r *ReqHandler) stat(req *transfer.StatReq) (*transfer.StatResp, error) {
	df, err := model.GetDataFile(req.DataFileID, r.db.session)
	switch {
	case err != nil:
		return nil, fmt.Errorf("Unknown id %s", req.DataFileID)
	case !ownerGaveAccessTo(df.Owner, r.user, r.db.session):
		return nil, fmt.Errorf("You do not have permission to access this datafile %s", req.DataFileID)
	default:
		return r.respStat(df), nil
	}
}

func (r *ReqHandler) respStat(df *model.DataFile) *transfer.StatResp {
	return &transfer.StatResp{
		DataFileID: df.Id,
		Name:       df.Name,
		DataDirs:   df.DataDirs,
		Checksum:   df.Checksum,
		Size:       df.Size,
		Birthtime:  df.Birthtime,
		MTime:      df.MTime,
	}
}
