package database

import "fmt"

func (s *Service) GetSystemLogs() []SystemLog {
	var logs []SystemLog
	s.orm.Order("ID desc").Find(&logs)
	return logs
}

func (s *Service) GetJobStats() *JobStats {
	stats := new(JobStats)
	err := s.orm.Model(&Run{}).
		Select(`COUNT(DISTINCT runs.id) AS total_runs,
				COUNT(logs.id) AS total_logs,
				SUM(CASE WHEN logs.log_severity_id = 2 THEN 1 ELSE 0 END) AS warning_runs,
				SUM(CASE WHEN logs.log_severity_id = 3 THEN 1 ELSE 0 END) AS error_runs,
				SUM(CASE WHEN logs.log_type_id = 1 THEN 1 ELSE 0 END) AS general_runs,
				SUM(CASE WHEN logs.log_type_id = 2 THEN 1 ELSE 0 END) AS restic_runs,
				SUM(CASE WHEN logs.log_type_id = 3 THEN 1 ELSE 0 END) AS custom_runs,
				SUM(CASE WHEN logs.log_type_id = 4 THEN 1 ELSE 0 END) AS prune_runs,
				SUM(CASE WHEN logs.log_type_id = 5 THEN 1 ELSE 0 END) AS check_runs`).
		Joins("LEFT JOIN logs ON runs.id = logs.run_id").
		Scan(&stats).Error
	if err != nil {
		fmt.Print(err)
	}
	return stats
}
