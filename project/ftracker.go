package project

import (
	"crypto/md5"
	"encoding/json"
	"github.com/materials-commons/gohandy/handyfile"
	"github.com/materials-commons/materials"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"path/filepath"
	"time"
)

type projectFileInfo struct {
	Path     string
	Size     int64
	Checksum []byte
	ModTime  time.Time
}

func WalkProject(project materials.Project) {
	db, _ := leveldb.OpenFile("/tmp/project.db", nil)
	defer db.Close()
	hasher := md5.New()
	filepath.Walk(project.Path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			key := hasher.Sum([]byte(path))
			checksum, _ := handyfile.Hash(hasher, path)
			pinfo := &projectFileInfo{
				Path:     path,
				Size:     info.Size(),
				Checksum: checksum,
				ModTime:  info.ModTime(),
			}
			value, _ := json.Marshal(pinfo)
			db.Put(key, value, nil)
		}
		return nil
	})
}
