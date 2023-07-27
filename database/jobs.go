package database

import (
	"time"

	"gorm.io/gorm"
)

func (s *Service) GetJob(id uint64) *Job {
	var job Job
	JobQuery(s.orm, &job, id)
	return &job
}

func (s *Service) DeleteJob(id uint64) {
	s.orm.Delete(&Job{}, id)
}

func (s *Service) GetJobsSelect(jobSelect ...string) []Job {
	var jobs []Job
	JobsSelectQuery(s.orm, &jobs, jobSelect...)
	return jobs
}

func (s *Service) GetJobs() []Job {
	var jobs []Job
	JobQuery(s.orm, &jobs)
	return jobs
}

func JobsSelectQuery(orm *gorm.DB, jobs *[]Job, jobSelect ...string) {
	orm.Select(jobSelect).Find(&jobs)
}

func JobQuery(orm *gorm.DB, value interface{}, conds ...interface{}) {
	orm.Preload("RetentionPolicy").Preload("CompressionType").Preload(
		"PreCommands", func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", 1).Order("commands.sort_id")
		},
	).Preload(
		"PostCommands", func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", 2).Order("commands.sort_id")
		},
	).Preload(
		"Runs", "start_time > ?", time.Now().UnixMilli()-TimeToGoBackInMilliseconds, func(db *gorm.DB) *gorm.DB {
			return db.Order("runs.id DESC")
		},
	).Preload(
		"Runs.Logs",
	).Order("Description").Find(value, conds)
}
