package statsHandler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nodeboxhq/nodebox-dashboard/internal/services/stats"
)

func Stats(service *stats.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.DefaultQuery("limit", "60")
		limit, err := strconv.Atoi(query)

		if err != nil {
			limit = 60
		}

		stats, err := service.Stats(limit)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, stats)
	}
}
