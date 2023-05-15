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

// AddFile will add a new AddedFile to the added_files table. Note that the path column is unique
// and the add will fail if there is already an entry matching that path.
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

// GetFileByPath returns the AddFile matching path.
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

// RemoveAll will remove all added files in the add_files tables.
func (s *GormAddedFileStor) RemoveAll() error {
	return s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.AddedFile{}).Error
}

// RemoveByPath will remove an entry in the add_files table that matches path.
func (s *GormAddedFileStor) RemoveByPath(path string) error {
	return s.db.Where("path = ?", path).Delete(&model.AddedFile{}).Error
}
