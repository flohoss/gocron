package database

import "gorm.io/gorm/clause"

func (s *Service) GetJob(id uint64) *Job {
	var job Job
	s.orm.Limit(1).Preload(clause.Associations).Find(&job, id)
	return &job
}

func (s *Service) GetJobs() []Job {
	var jobs []Job
	s.orm.Preload(clause.Associations).Order("Description").Find(&jobs)
	return jobs
}

func (s *Service) CreateOrUpdateJob(job *Job) error {
	if err := s.orm.Save(job).Error; err != nil {
		return err
	}
	return nil
}
