package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
)

func (h *ReqHandler) stat(req *transfer.StatReq) (*transfer.StatResp, error) {
	df, err := model.GetDataFile(req.DataFileID, h.session)
	switch {
	case err != nil:
		return nil, fmt.Errorf("Unknown id %s", req.DataFileID)
	case !OwnerGaveAccessTo(df.Owner, h.user, h.session):
		return nil, fmt.Errorf("You do not have permission to access this datafile %s", req.DataFileID)
	default:
		return respStat(df), nil
	}
}

func respStat(df *model.DataFile) *transfer.StatResp {
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
