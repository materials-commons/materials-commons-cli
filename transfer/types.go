package transfer

type StartHeader struct {
	ProjectID string
	Owner string
	ApiKey string
}

type DataFile struct {
	Path string
	DataDirID string
	ID string
	Checksum string
	Size int64
}

type SendStart struct {
	Header StartHeader
	DataFile DataFile
}

type SendStartResponse struct {
	Offset int64
	DataFileID string
	Status error
}

type MoveStart struct {
	Header StartHeader
	DataFile DataFile
	ToDataDirID string
}

type MoveStartResponse struct {
	DataDirID string
	Status error
}

type SendEndResponse struct {
	Checksum string
	Size int64
	DataFileID string
	Status error
}
