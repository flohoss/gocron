package database

func (s *Service) GetSelectOptionsCompressionTypes() []SelectOption {
	var temp []SelectOption
	var ct []CompressionType
	s.orm.Order("ID").Find(&ct)
	for _, option := range ct {
		temp = append(temp, SelectOption{
			Name:  option.Compression,
			Value: option.ID,
		})
	}
	return temp
}
