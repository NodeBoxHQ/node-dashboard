package hosts

import (
	"github.com/gin-gonic/gin"
	"github.com/nodeboxhq/nodebox-dashboard/internal/services/host"
)

func HostInfo(service *host.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := service.HostInfo()
		c.JSON(200, host)
	}
}
