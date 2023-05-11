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

func (s *GormAddedFileStor) AddFile(path, reason string) (*model.AddedFile, error) {
	fileToAdd := &model.AddedFile{
		Path:   path,
		Reason: reason,
	}

	err := mcdb.WithTxRetryDefault(s.db, func(tx *gorm.DB) error {
		return tx.Create(fileToAdd).Error
	})

	return fileToAdd, err
}
