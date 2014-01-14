package materials

import (
	"fmt"
	//"os"
	//"path/filepath"
	//"strings"
	"github.com/materials-commons/gohandy/marshaling"
	"github.com/materials-commons/materials/transfer"
)

var _ = fmt.Println

type projectUploader struct {
	Project
	marshaling.MarshalUnmarshaler
}

func (p Project) Upload2(m marshaling.MarshalUnmarshaler) error {
	pupload := &projectUploader{
		Project:            p,
		MarshalUnmarshaler: m,
	}

	ids, err := pupload.createProject()
	var _ = ids
	if err != nil {
		return err
	}

	if true {
		return nil
	}

	/*
		dir2id := make(map[string]string)
		dir2id[p.Path] = ids.DatadirId
		filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if path != p.Path {

				}
			}
			return nil
		})
	*/
	return nil
}

func (p *projectUploader) createProject() (string, error) {
	createProjectRequest := transfer.CreateProjectReq{
		Name: p.Name,
	}
	req := transfer.Request{&createProjectRequest}
	err := p.Marshal(&req)
	if err != nil {
		return "", err
	}

	var resp transfer.Response

	err = p.Unmarshal(&resp)
	if err != nil {
		return "", nil
	}

	if resp.Type != transfer.ROk {

	}

	return "", nil
}
