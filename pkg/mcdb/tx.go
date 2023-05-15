package mcdb

import (
	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"gorm.io/gorm"
)

// WithTxRetryDefault attempts to execute a database transaction up to config.GetTxRetry times
// before failing it.
func WithTxRetryDefault(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return WithTxRetry(db, config.GetTxRetry(), fn)
}

// WithTxRetry attempts to execute a database transaction retryCount times. If retryCount < 3, then
// it is set to 3.
func WithTxRetry(db *gorm.DB, retryCount int, fn func(tx *gorm.DB) error) error {
	var err error

	if retryCount < 3 {
		retryCount = 3
	}

	for i := 0; i < retryCount; i++ {
		err = db.Transaction(fn)
		if err == nil {
			break
		}
	}

	return err
}
