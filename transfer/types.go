package transfer

type TransferType int

const (
	Upload TransferType = iota
	Download
	Move
	Delete
)

type StartHeader struct {
	ProjectID string
	Owner     string
	ApiKey    string
}

type FileBlock struct {
	DataFileID string
	Bytes      []byte
}

type DataFile struct {
	Path      string
	DataDirID string
	ID        string
	Checksum  string
	Size      int64
}

type Command struct {
	Header   StartHeader
	DataFile DataFile
	Type     TransferType
}

type SendStartResponse struct {
	Offset     int64
	DataFileID string
	Status     error
}

type MoveStart struct {
	Header      StartHeader
	DataFile    DataFile
	ToDataDirID string
}

type MoveStartResponse struct {
	DataDirID string
	Status    error
}

type SendEndResponse struct {
	Checksum   string
	Size       int64
	DataFileID string
	Status     error
}
