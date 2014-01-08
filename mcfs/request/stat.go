package request

import (
	"fmt"
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
)

func (r *ReqHandler) stat(req transfer.Request) ReqStateFN {
	switch t := req.Req.(type) {
	case transfer.StatReq:
		df, err := model.GetDataFile(t.DataFileID, r.db.session)
		if err != nil {
			return r.badRequestNext(fmt.Errorf("Unknown id %s", t.DataFileID))
		}
		r.respStat(df)
		return r.nextCommand()
	default:
		return r.badRequestNext(fmt.Errorf("Bad request data for type %s", req.Type))
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
