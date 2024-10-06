package services

import (
	"github.com/nodeboxhq/nodebox-dashboard/internal/db/models"
	"github.com/nodeboxhq/nodebox-dashboard/internal/logger"
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

		var hostInfo models.Host

		if hostService != nil {
			hostInfo = hostService.HostInfo()
			nodeName = hostInfo.Node
		}

		*s = node.Service{DB: db, NodeName: nodeName}

		if hostInfo.Hostname != "" {
			logger.L.Info().Msgf("Hostname: %s, Node: %s, IPv4 (Private): %s, IPv6 (Private): %s, IPv4 (Public): %s, IPv6 (Public): %s", hostInfo.Hostname, hostInfo.Node, hostInfo.PrivateIPv4, hostInfo.PrivateIPv6, hostInfo.IPv4, hostInfo.IPv6)
		}
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
