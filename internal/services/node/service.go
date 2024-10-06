package node

import (
	"gorm.io/gorm"
)

type Service struct {
	DB       *gorm.DB
	NodeName string
}

func (s *Service) GetNodeInfo() (interface{}, error) {
	if s.NodeName == "linea" || s.NodeName == "pc" {
		return LineaInfo(), nil
	} else if s.NodeName == "dusk" {
		return DuskInfo(), nil
	} else if s.NodeName == "juneo" {
		return JuneoInfo(), nil
	}

	return nil, nil
}
