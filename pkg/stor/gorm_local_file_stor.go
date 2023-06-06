package stor

import (
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"gorm.io/gorm"
)

type GormLocalFileStor struct {
	db *gorm.DB
}

func NewGormLocalFileStor(db *gorm.DB) *GormLocalFileStor {
	return &GormLocalFileStor{db: db}
}

func (s *GormLocalFileStor) ListPaged(fn func(f *model.LocalFile) error) error {
	return listPaged(s.db, fn)
}

func (s *GormLocalFileStor) GetRemoteFileByPath(path string) (*model.LocalFile, error) {
	var f model.LocalFile
	err := s.db.Where("path = ?", path).First(&f).Error
	return &f, err
}
