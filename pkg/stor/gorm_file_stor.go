package stor

import (
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"gorm.io/gorm"
)

type GormFileStor struct {
	db *gorm.DB
}

func NewGormFileStor(db *gorm.DB) *GormFileStor {
	return &GormFileStor{db: db}
}

func (s *GormFileStor) GetFileByPath(path string) (*model.File, error) {
	var f model.File

	if err := s.db.Where("path = ?", path).First(&f).Error; err != nil {
		return nil, err
	}

	return &f, nil
}

func (s *GormFileStor) AddFile(f model.File) (*model.File, error) {
	fileToAdd := f
	err := mcdb.WithTxRetryDefault(s.db, func(tx *gorm.DB) error {
		return tx.Create(&fileToAdd).Error
	})
	return &fileToAdd, err
}
