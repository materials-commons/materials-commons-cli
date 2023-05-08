package model

type Conflict struct {
	ID       uint  `json:"id"`
	RemoteID uint  `json:"remote_id"`
	FileID   uint  `json:"file_id"`
	File     *File `json:"file" gorm:"foreignKey:FileID;references:ID"`
}
