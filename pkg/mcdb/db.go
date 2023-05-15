package mcdb

import (
	"log"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&model.File{}, &model.Conflict{}, &model.AddedFile{}, &model.IgnoredFile{})
}

// MustConnectToDB connects to the database and returns the db instance. If it cannot connect them
// it logs a fatal error and exists.
func MustConnectToDB() *gorm.DB {
	var (
		err error
		db  *gorm.DB
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	dsn := config.GetProjectDBPath()

	if db, err = gorm.Open(sqlite.Open(dsn), gormConfig); err != nil {
		log.Fatalf("Failed to open db (%s): %s", dsn, err)
	}

	return db
}
