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

type ProjectFileStatus int

const (
	Synchronized ProjectFileStatus = iota
	Unsynchronized
	New
	Deleted
)

func (pfs ProjectFileStatus) String() string {
	switch {
	case pfs&Synchronized == Synchronized:
		return "Synchronized"
	case pfs&Unsynchronized == Unsynchronized:
		return "Unsynchronized"
	case pfs&New == New:
		return "New"
	case pfs&Deleted == Deleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}

type ProjectFileLocation int

const (
	LocalOnly ProjectFileLocation = iota
	RemoteOnly
	LocalAndRemote
	LocalAndRemoteUnknown
)

func (pfl ProjectFileLocation) String() string {
	switch {
	case pfl&LocalOnly == LocalOnly:
		return "LocalOnly"
	case pfl&RemoteOnly == RemoteOnly:
		return "RemoteOnly"
	case pfl&LocalAndRemote == LocalAndRemote:
		return "LocalAndRemote"
	case pfl&LocalAndRemoteUnknown == LocalAndRemoteUnknown:
		return "LocalAndRemoteUnknown"
	default:
		return "Unknown"
	}
}

type ProjectFileInfo struct {
	Path     string
	Size     int64
	Hash     string
	ModTime  time.Time
	Id       string
	Status   ProjectFileStatus
	Location ProjectFileLocation
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
