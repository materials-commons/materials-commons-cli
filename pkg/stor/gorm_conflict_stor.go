package stor

import (
	"time"

	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
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
	return mcdb.WithTxRetryDefault(s.db, func(tx *gorm.DB) error {
		return tx.Model(&model.Conflict{}).
			Where("path = ?", path).
			Update("resolved_at", time.Now()).Error
	})
}

func (s *GormConflictStor) ResolveAllConflicts() error {
	return mcdb.WithTxRetryDefault(s.db, func(tx *gorm.DB) error {
		return tx.Model(&model.Conflict{}).
			Where("resolved_at IS NULL").
			Update("resolved_at", time.Now()).Error
	})
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

func (s *GormConflictStor) GetConflictByPath(path string) (*model.Conflict, error) {
	var f model.Conflict
	err := s.db.Where("path = ?", path).First(&f).Error
	return &f, err
}
