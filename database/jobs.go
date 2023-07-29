package database

import (
	"gorm.io/gorm"
)

func (s *Service) GetJobs() []Job {
	var jobs []Job
	JobBaseSelect(s.orm, &jobs)

	for i := range jobs {
		latestRun := jobs[i].getLatestRun(s.orm)
		jobs[i].Status = latestRun.getHighestLogSeverity(s.orm)
	}

	return jobs
}

func (s *Service) GetJob(id uint64) *Job {
	var job Job
	JobBaseSelect(s.orm, &job, id)
	return &job
}

func (s *Service) DeleteJob(id uint64) {
	s.orm.Delete(&Job{}, id)
}

func JobBaseSelect(orm *gorm.DB, value interface{}, conds ...interface{}) {
	orm.Preload(
		"PreCommands", func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", 1).Order("commands.sort_id")
		},
	).Preload(
		"PostCommands", func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", 2).Order("commands.sort_id")
		},
	).Order("description").Find(value, conds...)
}

func (j *Job) getLatestRun(orm *gorm.DB) Run {
	var latestRun Run
	orm.Where("job_id = ?", j.ID).Order("start_time DESC").Limit(1).Find(&latestRun)
	return latestRun
}

func (r *Run) getHighestLogSeverity(orm *gorm.DB) LogSeverity {
	var highestSeverity LogSeverity
	orm.Model(&Log{}).Select("MAX(log_severity)").Where("run_id = ?", r.ID).Row().Scan(&highestSeverity)
	return highestSeverity
}
