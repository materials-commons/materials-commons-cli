package stor

import (
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"gorm.io/gorm"
)

type GormRemoteFileStor struct {
	db *gorm.DB
}

func NewGormRemoteFileStor(db *gorm.DB) *GormRemoteFileStor {
	return &GormRemoteFileStor{db: db}
}

func (s *GormRemoteFileStor) ListPaged(fn func(f *model.RemoteFile) error) error {
	return listPaged(s.db, fn)
}

func (s *GormRemoteFileStor) GetRemoteFileByPath(path string) (*model.RemoteFile, error) {
	var f model.RemoteFile

	err := s.db.Where("path = ?", path).First(&f).Error
	return &f, err
}
