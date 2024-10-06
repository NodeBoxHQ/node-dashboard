package services

import (
	"github.com/nodeboxhq/nodebox-dashboard/internal/services/host"
	"github.com/nodeboxhq/nodebox-dashboard/internal/services/node"
	"github.com/nodeboxhq/nodebox-dashboard/internal/services/stats"
	"gorm.io/gorm"
)

type ServiceRegistry struct {
	StatsService *stats.Service
	HostService  *host.Service
	NodeService  *node.Service
}

func NewService[T any](db *gorm.DB, dependencies ...interface{}) *T {
	var service T

	switch s := any(&service).(type) {
	case *stats.Service:
		*s = stats.Service{DB: db}
	case *host.Service:
		*s = host.Service{DB: db}
	case *node.Service:
		var hostService *host.Service
		for _, dep := range dependencies {
			if hs, ok := dep.(*host.Service); ok {
				hostService = hs
				break
			}
		}

		nodeName := ""
		if hostService != nil {
			hostInfo := hostService.HostInfo()
			nodeName = hostInfo.Node
		}

		*s = node.Service{DB: db, NodeName: nodeName}
	}

	return &service
}

func NewServiceRegistry(db *gorm.DB) *ServiceRegistry {
	hostService := NewService[host.Service](db)

	return &ServiceRegistry{
		StatsService: NewService[stats.Service](db),
		HostService:  hostService,
		NodeService:  NewService[node.Service](db, hostService),
	}
}
