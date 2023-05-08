package mcdb

import (
	"log"
	"path/filepath"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&model.File{}, &model.Conflict{}, &model.AddedFile{}, &model.IgnoredFile{})
}

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

func MustConnectToDBAtPath(dir string) *gorm.DB {
	var (
		err error
		db  *gorm.DB
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	dsn := filepath.Join(dir, "repository.db")

	if db, err = gorm.Open(sqlite.Open(dsn), gormConfig); err != nil {
		log.Fatalf("Failed to open db (%s): %s", dsn, err)
	}

	return db
}
