package node

import (
	"gorm.io/gorm"
)

type Service struct {
	DB       *gorm.DB
	NodeName string
}

func (s *Service) GetNodeInfo() (interface{}, error) {
	if s.NodeName == "linea" {
		return LineaInfo(), nil
	} else if s.NodeName == "dusk" {
		return DuskInfo(), nil
	} else if s.NodeName == "juneo" {
		return JuneoInfo(), nil
	} else if s.NodeName == "hyperliquid" || s.NodeName == "pc" {
		return HyperliquidInfo()
	}

	return nil, nil
}
