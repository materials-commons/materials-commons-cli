package stor

import (
	"gorm.io/gorm"
)

// listPaged implements database paging in a generic method. Page size is always 100. The method
// will stop if the callback method returns a non-nil error. This gives the user control over
// when the method should stop.
func listPaged[M Model](db *gorm.DB, fn func(m *M) error) error {
	var items []M
	offset := 0
	pageSize := 100
	for {
		if err := db.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
			return err
		}

		if len(items) == 0 {
			break
		}

		for _, f := range items {
			if err := fn(&f); err != nil {
				break
			}
		}
		offset = offset + pageSize
	}

	return nil
}
