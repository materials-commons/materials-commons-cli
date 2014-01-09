package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
)

func (r *ReqHandler) stat(req *transfer.StatReq) ReqStateFN {
	df, err := model.GetDataFile(req.DataFileID, r.db.session)
	switch {
	case err != nil:
		return r.badRequestNext(fmt.Errorf("Unknown id %s", req.DataFileID))
	case !ownerGaveAccessTo(df.Owner, r.user, r.db.session):
		return r.badRequestNext(fmt.Errorf("You do not have permission to access this datafile %s", req.DataFileID))
	default:
		r.respStat(df)
		return r.nextCommand()
	}
}

func (r *ReqHandler) respStat(df *model.DataFile) {
	statResp := &transfer.StatResp{
		DataFileID: df.Id,
		Name:       df.Name,
		DataDirs:   df.DataDirs,
		Checksum:   df.Checksum,
		Size:       df.Size,
		Birthtime:  df.Birthtime,
		MTime:      df.MTime,
	}
	r.respOk(statResp)
}
