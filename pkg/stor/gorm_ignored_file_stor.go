package stor

import (
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"gorm.io/gorm"
)

type GormIgnoredFileStor struct {
	db *gorm.DB
}

func NewGormIgnoredFileStor(db *gorm.DB) *GormIgnoredFileStor {
	return &GormIgnoredFileStor{db: db}
}

func (s *GormIgnoredFileStor) FileIsIgnored(path string) bool {
	var ignoredFile model.IgnoredFile

	if err := s.db.Where("path = ?", path).First(&ignoredFile).Error; err != nil {
		return false
	}

	return true
}
