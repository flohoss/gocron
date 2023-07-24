package database

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) GetJob(id uint64) *Job {
	var job Job
	s.orm.Limit(1).Preload(
		"PreCommands", func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", 1).Order("commands.sort_id")
		},
	).Preload(
		"PostCommands", func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", 2).Order("commands.sort_id")
		},
	).Preload(clause.Associations).Find(&job, id)
	return &job
}

func (s *Service) DeleteJob(id uint64) {
	s.orm.Delete(&Job{}, id)
}

func (s *Service) GetJobs() []Job {
	var jobs []Job
	sevenDaysAgo := time.Now().UnixMilli() - 604800000
	s.orm.Preload(
		"PreCommands", func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", 1).Order("commands.sort_id")
		},
	).Preload(
		"PostCommands", func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", 2).Order("commands.sort_id")
		},
	).Preload(
		"Runs", "start_time > ?", sevenDaysAgo, func(db *gorm.DB) *gorm.DB {
			return db.Order("runs.id DESC")
		},
	).Preload(
		"Runs.Logs",
	).Preload(
		"Runs.Logs.LogType",
	).Preload(
		"Runs.Logs.LogSeverity",
	).Order("Description").Find(&jobs)
	return jobs
}
