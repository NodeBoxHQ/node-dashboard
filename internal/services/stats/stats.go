package stats

import (
	"context"
	"time"

	"github.com/nodeboxhq/nodebox-dashboard/internal/db/models"
	"github.com/nodeboxhq/nodebox-dashboard/internal/logger"
	"github.com/nodeboxhq/nodebox-dashboard/internal/utils"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func (s *Service) Stats(limit int) ([]models.Stats, error) {
	var stats []models.Stats

	if limit == 0 {
		limit = 60
	} else if limit > 720 {
		limit = 720
	}

	err := s.DB.Raw(`
		SELECT * FROM (
			SELECT * FROM stats ORDER BY created_at DESC LIMIT ?
		) AS subquery
		ORDER BY created_at ASC
	`, limit).Scan(&stats).Error

	return stats, err
}

func (s *Service) CollectStats() {
	store := models.Stats{}

	c, err := cpu.Percent(time.Second, false)
	cpuAvg := 0.0
	for _, core := range c {
		cpuAvg += core
	}
	store.CPU = int(cpuAvg / float64(len(c)))
	if err != nil {
		logger.L.Err(err).Msg("Failed to collect CPU stats")
	}

	r, err := mem.VirtualMemory()
	if err != nil {
		logger.L.Err(err).Msg("Failed to collect memory stats")
	}
	store.Memory = int(r.UsedPercent)

	d, err := disk.Usage("/")
	if err != nil {
		logger.L.Err(err).Msg("Failed to collect disk stats")
	}
	store.Storage = int(d.UsedPercent)

	netIO, err := net.IOCounters(false)

	if err != nil {
		logger.L.Err(err).Msg("Failed to collect network stats")
	}

	var network int

	for _, io := range netIO {
		network += int(io.BytesSent + io.BytesRecv)
	}

	store.Network = network

	uptime, err := utils.GetUptimeSeconds()

	if err != nil {
		logger.L.Err(err).Msg("Failed to collect uptime")
	}

	store.Uptime = int(uptime)

	tx := s.DB.Begin()

	if err := tx.Create(&store).Error; err != nil {
		tx.Rollback()
		logger.L.Err(err).Msg("Failed to save stats")
		return
	}

	var count int64
	if err := tx.Model(&models.Stats{}).Count(&count).Error; err != nil {
		tx.Rollback()
		logger.L.Err(err).Msg("Failed to count stats")
		return
	}

	if count > 720 {
		deleteCount := count - 720
		var oldestIDs []uint
		if err := tx.Model(&models.Stats{}).Order("created_at asc").Limit(int(deleteCount)).Pluck("id", &oldestIDs).Error; err != nil {
			tx.Rollback()
			logger.L.Err(err).Msg("Failed to fetch oldest stats")
			return
		}
		if err := tx.Delete(&models.Stats{}, oldestIDs).Error; err != nil {
			tx.Rollback()
			logger.L.Err(err).Msg("Failed to delete old stats")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.L.Err(err).Msg("Failed to commit transaction")
		return
	}

	logger.L.Debug().Msgf("Stats collected: %+v", store)
}

func (s *Service) StartStatsCollection(ctx context.Context) {
	s.CollectStats()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.CollectStats()
		case <-ctx.Done():
			logger.L.Info().Msg("Stopping stats collection")
			return
		}
	}
}
