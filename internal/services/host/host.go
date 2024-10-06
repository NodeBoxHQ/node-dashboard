package host

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nodeboxhq/nodebox-dashboard/internal/db/models"
	"github.com/nodeboxhq/nodebox-dashboard/internal/logger"
	"github.com/nodeboxhq/nodebox-dashboard/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func (s *Service) HostInfo() models.Host {
	var host models.Host
	s.DB.First(&host, 1)
	return host
}

func (s *Service) CollectHostInfo() error {
	store := models.Host{}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}
	store.Hostname = hostname

	owner := "Unknown"
	if strings.Contains(hostname, "-") {
		split := strings.Split(hostname, "-")
		owner = strings.Join(split[:len(split)-1], "-")
	}
	store.Owner = owner

	ipv4, ipv6, err := utils.GetPublicIPs()
	if err != nil {
		ipv4 = "Unknown"
		ipv6 = "Unknown"
	}
	store.IPv4 = ipv4
	store.IPv6 = ipv6

	privateIpv4, privateIpv6 := utils.GetPrivateIPs()
	store.PrivateIPv4 = privateIpv4
	store.PrivateIPv6 = privateIpv6

	node := "Unknown"
	if strings.Contains(hostname, "-") {
		split := strings.Split(hostname, "-")
		node = split[len(split)-1]
	}
	store.Node = node

	var existingHost models.Host
	if err := s.DB.First(&existingHost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.DB.Create(&store).Error; err != nil {
				return fmt.Errorf("failed to create host info: %w", err)
			}
		} else {
			return fmt.Errorf("failed to query existing host info: %w", err)
		}
	} else {
		if err := s.DB.Model(&existingHost).Updates(store).Error; err != nil {
			return fmt.Errorf("failed to update host info: %w", err)
		}
	}

	return nil
}

func (s *Service) StartHostInfoCollection(ctx context.Context) {
	s.CollectHostInfo()

	ticker := time.NewTicker(120 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.CollectHostInfo()
		case <-ctx.Done():
			logger.L.Info().Msg("Stopping host info collection")
			return
		}
	}
}
