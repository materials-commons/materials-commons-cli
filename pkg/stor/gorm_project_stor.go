package stor

import (
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"gorm.io/gorm"
)

type GormProjectStor struct {
	db *gorm.DB
}

func NewGormProjectStor(db *gorm.DB) *GormProjectStor {
	return &GormProjectStor{db: db}
}

func (s *GormProjectStor) GetProject() (*model.Project, error) {
	var p model.Project
	err := s.db.First(&p).Error
	return &p, err
}
