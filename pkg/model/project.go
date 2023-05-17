package model

type Project struct {
	ID   uint   `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (Project) TableName() string {
	return "project"
}
