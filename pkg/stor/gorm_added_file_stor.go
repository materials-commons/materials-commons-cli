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

// ListPaged pages through all the add files, the callback method is called on each
// conflict. If the method returns a non-nil error then ListPaged will immediately
// stop execution.
func (s *GormAddedFileStor) ListPaged(fn func(f *model.AddedFile) error) error {
	return listPaged(s.db, fn)
}
