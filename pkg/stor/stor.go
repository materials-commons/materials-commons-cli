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
	AddFile(path, reason, ftype string) (*model.AddedFile, error)
	GetFileByPath(path string) (*model.AddedFile, error)
	ListPaged(fn func(f *model.AddedFile) error) error
	RemoveAll() error
	RemoveByPath(path string) error
}

type ConflictFileStor interface {
	ResolveConflictByPath(path string) error
	ResolveAllConflicts() error
	ListPaged(fn func(conflict *model.Conflict) error) error
	GetConflictByPath(path string) (*model.Conflict, error)
}

type ProjectStor interface {
	GetProject() (*model.Project, error)
}

type RemoteStor interface {
	GetDefaultRemote() (*model.Remote, error)
	GetRemoteByUserServerUrl(email, serverUrl string) (*model.Remote, error)
	ListPaged(fn func(remote *model.Remote) error) error
}

type Model interface {
	model.Conflict | model.File | model.AddedFile | model.IgnoredFile
}
