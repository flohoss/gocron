package database

import (
	_ "embed"
	"os"

	"gitlab.unjx.de/flohoss/gobackup/internal/notify"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const Storage = "storage/"

func init() {
	os.Mkdir(Storage, os.ModePerm)
}

type Service struct {
	orm           *gorm.DB
	notifyService notify.Notify
	identifier    string
}

func MigrateDatabase(notifyService notify.Notify, identifier string) (*Service, error) {
	db, err := gorm.Open(sqlite.Open(Storage+"db.sqlite3?_pragma=foreign_keys(1)"), &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&SystemLog{})
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Run{})
	db.AutoMigrate(&Command{})
	db.AutoMigrate(&Job{})
	return &Service{orm: db, notifyService: notifyService, identifier: identifier}, nil
}
