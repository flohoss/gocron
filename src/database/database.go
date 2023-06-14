package database

import (
	"os"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

const Storage = "storage/"

func init() {
	os.Mkdir(Storage, os.ModePerm)
}

func NewDatabaseConnection(location string) *gorm.DB {
	logger := zapgorm2.New(zap.L())
	db, err := gorm.Open(sqlite.Open(Storage+location+"?_pragma=foreign_keys(1)"), &gorm.Config{Logger: logger, SkipDefaultTransaction: true, PrepareStmt: true})
	if err != nil {
		zap.S().Fatal("Cannot connect to database", zap.Error(err))
	}
	return db
}
