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

type SelectOption struct {
	Value uint64 `json:"value"`
	Name  string `json:"name"`
}

//go:embed sql/inserts.sql
var inserts string

func MigrateDatabase() (*Service, error) {
	db, err := gorm.Open(sqlite.Open(Storage+"db.sqlite3?_pragma=foreign_keys(1)"), &gorm.Config{Logger: zapgorm2.New(zap.L()), SkipDefaultTransaction: true, PrepareStmt: true})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&RetentionPolicy{})
	db.AutoMigrate(&CompressionType{})
	db.AutoMigrate(&PreCommand{})
	db.AutoMigrate(&PostCommand{})
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&LogType{})
	db.AutoMigrate(&Run{})
	db.AutoMigrate(&Job{})

	if err := db.Exec(inserts).Error; err != nil {
		return nil, err
	}
	return &Service{orm: db}, nil
}
