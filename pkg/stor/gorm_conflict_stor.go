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

// ListPaged pages through all the conflicts, the callback method is called on each
// conflict. If the method returns a non-nil error then ListPaged will immediately
// stop execution.
func (s *GormConflictStor) ListPaged(fn func(conflict *model.Conflict) error) error {
	return listPaged(s.db, fn)
}

func (s *GormConflictStor) GetConflictByPath(path string) (*model.Conflict, error) {
	var f model.Conflict
	err := s.db.Where("path = ?", path).First(&f).Error
	return &f, err

}
