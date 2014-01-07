package transfer

import (
	"encoding/gob"
	"fmt"
	"time"
)

func init() {
	gob.Register(Request{})
	gob.Register(Response{})

	gob.Register(UploadReq{})
	gob.Register(UploadResp{})

	gob.Register(DownloadReq{})
	gob.Register(DownloadResp{})

	gob.Register(MoveReq{})
	gob.Register(MoveResp{})

	gob.Register(DeleteReq{})
	gob.Register(DeleteResp{})

	gob.Register(SendReq{})
	gob.Register(SendResp{})

	gob.Register(StatReq{})
	gob.Register(StatResp{})

	gob.Register(EndReq{})
	gob.Register(EndResp{})

	gob.Register(CreateFileReq{})
	gob.Register(CreateDirReq{})
	gob.Register(CreateProjectReq{})

	gob.Register(CreateResp{})

	gob.Register(LoginReq{})

	gob.Register(StartResp{})
}

type ResponseType int

const (
	ROk ResponseType = iota
	RError
	RFatal
)

type RequestType int

const (
	Upload RequestType = iota
	RestartUpload
	Download
	Move
	Delete
	Send
	Stat
	Login
	Done
	CreateFile
	CreateDir
	CreateProject
	Error
	Logout
	Close
)

var types = map[RequestType]bool{
	Upload:        true,
	RestartUpload: true,
	Download:      true,
	Move:          true,
	Delete:        true,
	Send:          true,
	Stat:          true,
	Login:         true,
	Done:          true,
	CreateFile:    true,
	CreateDir:     true,
	CreateProject: true,
	Error:         true,
	Logout:        true,
	Close:         true,
}

var requestTypeString = map[RequestType]string{
	Upload:        "Upload",
	RestartUpload: "RestartUpload",
	Download:      "Download",
	Move:          "Move",
	Delete:        "Delete",
	Send:          "Send",
	Stat:          "Stat",
	Login:         "Login",
	Done:          "Done",
	CreateFile:    "CreateFile",
	CreateDir:     "CreateDir",
	CreateProject: "CreateProject",
	Error:         "Error",
	Logout:        "Logout",
	Close:         "Close",
}

func (t RequestType) String() string {
	s := requestTypeString[t]
	if s == "" {
		return "Unknown"
	}

	return s
}

func ValidType(t RequestType) bool {
	return types[t]
}

type ItemType int

const (
	DataDir ItemType = iota
	DataFile
	Project
	DataSet
)

var itemTypes = map[ItemType]bool{
	DataDir:  true,
	DataFile: true,
	Project:  true,
	DataSet:  true,
}

func ValidItemType(t ItemType) bool {
	return itemTypes[t]
}

type Request struct {
	Type RequestType
	Req  interface{}
}

type Response struct {
	Type   ResponseType
	Status string
	Resp   interface{}
}

type UploadReq struct {
	DataFileID string
	Checksum   string
	Size       int64
}

type UploadResp struct {
	DataFileID string
}

type DownloadReq struct {
	Type ItemType
	ID   string
}

type DownloadResp struct {
	Ok bool
}

type MoveReq struct {
}

type MoveResp struct {
}

type DeleteReq struct {
}

type DeleteResp struct {
}

type SendReq struct {
	DataFileID string
	Bytes      []byte
}

type SendResp struct {
	BytesWritten int64
}

type StatReq struct {
	Path       string
	DataDirID  string
	DataFileID string
	Checksum   string
	Size       int64
}

type StatResp struct {
	DataFileID string
	Checksum   string
	Size       int64
	Birthtime  time.Time
	MTime      time.Time
}

type EndReq struct {
}

type EndResp struct {
	Ok bool
}

type CreateFileReq struct {
	ProjectID string
	DataDirID string
	Path      string
	Checksum  string
	Size      int64
}

type CreateDirReq struct {
	ProjectID string
	Path      string
}

type CreateProjectReq struct {
	Path string
}

type CreateResp struct {
	ID string
}

type LoginReq struct {
	ProjectID string
	User      string
	ApiKey    string
}

type StartResp struct {
	Ok bool
}
