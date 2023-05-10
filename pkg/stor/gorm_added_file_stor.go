package stor

import (
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"gorm.io/gorm"
)

type GormAddedFileStor struct {
	db *gorm.DB
}

func NewGormAddedFileStor(db *gorm.DB) *GormAddedFileStor {
	return &GormAddedFileStor{db: db}
}

func (s *GormAddedFileStor) AddFile(path string) (*model.AddedFile, error) {
	fileToAdd := &model.AddedFile{
		Path: path,
	}

	err := mcdb.WithTxRetryDefault(s.db, func(tx *gorm.DB) error {
		return tx.Create(fileToAdd).Error
	})

	return fileToAdd, err
}
