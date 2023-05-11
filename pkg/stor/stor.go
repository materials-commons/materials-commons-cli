package stor

import (
	"github.com/materials-commons/materials-commons-cli/pkg/model"
)

type IgnoredFileStor interface {
	FileIsIgnored(path string) bool
}

type FileStor interface {
	GetFileByPath(path string) (*model.File, error)
}

type AddedFileStor interface {
	AddFile(path, reason string) (*model.AddedFile, error)
}

type ConflictFileStor interface {
}
