package database

func (s *Service) DeleteCommand(id uint64) {
	s.orm.Delete(&Command{}, id)
}
