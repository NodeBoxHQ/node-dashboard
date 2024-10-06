package node

import (
	"github.com/gin-gonic/gin"
	"github.com/nodeboxhq/nodebox-dashboard/internal/services/node"
)

func NodeInfo(service *node.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := service.GetNodeInfo()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, info)
	}
}
