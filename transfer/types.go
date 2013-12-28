package transfer

type Type int

const (
	Upload Type = iota
	Download
	Move
	Delete
)

var commandTypes = map[Type]bool{
	Upload:   true,
	Download: true,
	Move:     true,
	Delete:   true,
}

func ValidType(t Type) bool {
	return commandTypes[t]
}

type StartHeader struct {
	ProjectID string
	User      string
	ApiKey    string
}

type FileBlock struct {
	DataFileID string
	Bytes      []byte
	Done       bool
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
	Type     Type
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
