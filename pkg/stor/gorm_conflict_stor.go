package stor

import (
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"gorm.io/gorm"
)

type GormConflictStor struct {
	db *gorm.DB
}

func NewGormConflictStor(db *gorm.DB) *GormConflictStor {
	return &GormConflictStor{db: db}
}

func (s *GormConflictStor) ResolveConflictByPath(path string) error {
	return nil
}

func (s *GormConflictStor) ResolveAllConflicts() error {
	return nil
}

func (s *GormConflictStor) ListPaged(fn func(conflict *model.Conflict) error) error {
	var conflictFiles []model.Conflict
	offset := 0
	pageSize := 100
	for {
		if err := s.db.Offset(offset).Limit(pageSize).Find(&conflictFiles).Error; err != nil {
			return err
		}

		if len(conflictFiles) == 0 {
			break
		}

		for _, f := range conflictFiles {
			if err := fn(&f); err != nil {
				break
			}
		}
		offset = offset + pageSize
	}

	return nil
}
