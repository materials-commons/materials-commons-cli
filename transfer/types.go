package transfer

import (
	"time"
)

type ResponseType int

const (
	RError ResponseType = iota
	RContinue
)

type RequestType int

const (
	Upload RequestType = iota
	Download
	Move
	Delete
	Send
	Stat
	Login
	Done
	Create
	Error
	Logout
)

var types = map[RequestType]bool{
	Upload:   true,
	Download: true,
	Move:     true,
	Delete:   true,
	Send:     true,
	Stat:     true,
	Login:    true,
	Done:     true,
	Create:   true,
	Error:    true,
	Logout:   true,
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
	Status error
	Resp   interface{}
}

type UploadReq struct {
}

type UploadResp struct {
	Offset     int64
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

type CreateReq struct {
	ProjectID string
	DataDirID string
	Path      string
	IsDir     bool
	Checksum  string
	Size      int64
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
