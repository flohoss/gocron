package database

import (
	"gorm.io/gorm/clause"
)

func (s *Service) GetJob(id uint64) *Job {
	var job Job
	s.orm.Limit(1).Preload(clause.Associations).Find(&job, id)
	return &job
}

func (s *Service) DeleteJob(id uint64) {
	s.orm.Delete(&Job{}, id)
}

func (s *Service) GetJobs() []Job {
	var jobs []Job
	s.orm.Preload(clause.Associations).Order("Description").Find(&jobs)
	return jobs
}
