package model

type LocalFile struct {
	ID         int
	ParentID   int
	Name       string
	Path       string
	ParentPath string
	MTime      float64 `gorm:"column:mtime"`
	Size       int
	Checksum   string
}

func (LocalFile) TableName() string {
	return "localtree"
}
