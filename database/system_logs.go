package database

func (s *Service) GetSystemLogs() []SystemLog {
	var logs []SystemLog
	s.orm.Order("ID desc").Find(&logs)
	return logs
}
