package materials

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/materials-commons/gohandy/file"
	"os"
	"path/filepath"
	"time"
)

var (
	BadProjectFileStatusString   = fmt.Errorf("Unknown string value for ProjectFileStatus")
	BadProjectFileLocationString = fmt.Errorf("Unknown string value for ProjectFileLocation")
)

type ProjectFileStatus int

const UnsetOption = 0

const (
	Synchronized ProjectFileStatus = iota
	Unsynchronized
	New
	Deleted
	UnknownFileStatus
)

var pfs2Strings = map[ProjectFileStatus]string{
	Synchronized:      "Synchronized",
	Unsynchronized:    "Unsynchronized",
	New:               "New",
	Deleted:           "Deleted",
	UnknownFileStatus: "UnknownFileStatus",
}

var pfsString2Value = map[string]ProjectFileStatus{
	"Synchronized":      Synchronized,
	"Unsynchronized":    Unsynchronized,
	"New":               New,
	"Deleted":           Deleted,
	"UnknownFileStatus": UnknownFileStatus,
}

func (pfs ProjectFileStatus) String() string {
	str, found := pfs2Strings[pfs]
	switch found {
	case true:
		return str
	default:
		return "Unknown"
	}
}

func String2ProjectFileStatus(pfs string) (ProjectFileStatus, error) {
	val, found := pfsString2Value[pfs]
	switch found {
	case true:
		return val, nil
	default:
		return -1, BadProjectFileStatusString
	}
}

type ProjectFileLocation int

const (
	LocalOnly ProjectFileLocation = iota
	RemoteOnly
	LocalAndRemote
	LocalAndRemoteUnknown
)

var pfl2Strings = map[ProjectFileLocation]string{
	LocalOnly:             "LocalOnly",
	RemoteOnly:            "RemoteOnly",
	LocalAndRemote:        "LocalAndRemote",
	LocalAndRemoteUnknown: "LocalAndRemoteUnknown",
}

var pflString2Value = map[string]ProjectFileLocation{
	"LocalOnly":             LocalOnly,
	"RemoteOnly":            RemoteOnly,
	"LocalAndRemote":        LocalAndRemote,
	"LocalAndRemoteUnknown": LocalAndRemoteUnknown,
}

func (pfl ProjectFileLocation) String() string {
	str, found := pfl2Strings[pfl]
	switch found {
	case true:
		return str
	default:
		return "Unknown"
	}
}

func String2ProjectFileLocation(pfl string) (p ProjectFileLocation, err error) {
	val, found := pflString2Value[pfl]
	switch found {
	case true:
		return val, nil
	default:
		return -1, BadProjectFileLocationString
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

type TrackingOptions struct {
	FileStatus   ProjectFileStatus
	FileLocation ProjectFileLocation
}

func (project *Project) Walk(options *TrackingOptions) error {
	fileStatus := Unsynchronized
	fileLocation := LocalOnly

	if options != nil {
		if options.FileStatus != UnsetOption {
			fileStatus = options.FileStatus
		}

		if options.FileLocation != UnsetOption {
			fileLocation = options.FileLocation
		}
	}

	filepath.Walk(project.Path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			key := md5.New().Sum([]byte(path))
			checksum, _ := file.Hash(md5.New(), path)
			pinfo := &ProjectFileInfo{
				Path:     path,
				Size:     info.Size(),
				Hash:     fmt.Sprintf("%x", checksum),
				ModTime:  info.ModTime(),
				Status:   fileStatus,
				Location: fileLocation,
			}
			value, _ := json.Marshal(pinfo)
			project.Put(key, value, nil)
		}
		return nil
	})

	return nil
}
