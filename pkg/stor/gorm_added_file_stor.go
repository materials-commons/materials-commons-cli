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

func (s *GormAddedFileStor) AddFile(path, reason, ftype string) (*model.AddedFile, error) {
	fileToAdd := &model.AddedFile{
		Path:   path,
		Reason: reason,
		FType:  ftype,
	}

	err := mcdb.WithTxRetryDefault(s.db, func(tx *gorm.DB) error {
		return tx.Create(fileToAdd).Error
	})

	return fileToAdd, err
}

func (s *GormAddedFileStor) GetFileByPath(path string) (*model.AddedFile, error) {
	var f model.AddedFile

	err := s.db.Where("path = ?", path).First(&f).Error
	return &f, err
}

func (s *GormAddedFileStor) ListPaged(fn func(f *model.AddedFile) error) error {
	var addedFiles []model.AddedFile
	offset := 0
	pageSize := 100
	for {
		if err := s.db.Offset(offset).Limit(pageSize).Find(&addedFiles).Error; err != nil {
			return err
		}

		if len(addedFiles) == 0 {
			break
		}

		for _, f := range addedFiles {
			if err := fn(&f); err != nil {
				break
			}
		}
		offset = offset + pageSize
	}

	return nil
}
