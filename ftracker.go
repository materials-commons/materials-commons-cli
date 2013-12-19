package materials

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/materials-commons/gohandy/handyfile"
	"os"
	"path/filepath"
	"time"
)

type ProjectFileInfo struct {
	Path    string
	Size    int64
	Hash    string
	ModTime time.Time
	Id      string
}

func (project *Project) Walk() error {
	filepath.Walk(project.Path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			key := md5.New().Sum([]byte(path))
			checksum, _ := handyfile.Hash(md5.New(), path)
			pinfo := &ProjectFileInfo{
				Path:    path,
				Size:    info.Size(),
				Hash:    fmt.Sprintf("%x", checksum),
				ModTime: info.ModTime(),
			}
			value, _ := json.Marshal(pinfo)
			project.Put(key, value, nil)
		}
		return nil
	})

	return nil
}
