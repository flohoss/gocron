package database

import (
	_ "embed"
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

type Service struct {
	orm *gorm.DB
}

func MigrateDatabase() (*Service, error) {
	db, err := gorm.Open(sqlite.Open(Storage+"db.sqlite3?_pragma=foreign_keys(1)"), &gorm.Config{Logger: zapgorm2.New(zap.L()), PrepareStmt: true})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&SystemLog{})
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Run{})
	db.AutoMigrate(&Command{})
	db.AutoMigrate(&Job{})
	return &Service{orm: db}, nil
}
