package database

import "strconv"

func (s *Service) GetSelectOptionsRetentionPolicies() []SelectOption {
	var temp []SelectOption
	var rp []RetentionPolicy
	s.orm.Order("ID").Find(&rp)
	for _, option := range rp {
		temp = append(temp, SelectOption{
			Name:  option.Description,
			Value: strconv.FormatUint(uint64(option.ID), 10),
		})
	}
	return temp
}
