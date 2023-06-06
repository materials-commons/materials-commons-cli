package model

type RemoteFile struct {
	ID         int
	ParentID   int
	Name       string
	Path       string
	ParentPath string
	MTime      float64 `gorm:"column:mtime"`
	Size       int
	Checksum   string
	OType      string `gorm:"column:otype"`
}

func (RemoteFile) TableName() string {
	return "remotetree"
}
