package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	hostsHandler "github.com/nodeboxhq/nodebox-dashboard/internal/handlers/hosts"
	nodeHandler "github.com/nodeboxhq/nodebox-dashboard/internal/handlers/node"
	statsHandler "github.com/nodeboxhq/nodebox-dashboard/internal/handlers/stats"
	"github.com/nodeboxhq/nodebox-dashboard/internal/services/host"
	nodeService "github.com/nodeboxhq/nodebox-dashboard/internal/services/node"
	statsService "github.com/nodeboxhq/nodebox-dashboard/internal/services/stats"
)

var EmbeddedWebFS embed.FS

func RegisterRoutes(r *gin.Engine, environment string, statsService *statsService.Service,
	hostService *host.Service, nodeService *nodeService.Service) {
	apiGroup := r.Group("/api")
	apiGroup.GET("/stats", statsHandler.Stats(statsService))
	apiGroup.GET("/host", hostsHandler.HostInfo(hostService))
	apiGroup.GET("/node", nodeHandler.NodeInfo(nodeService))

	if environment == "dev" || environment == "development" {
		r.NoRoute(func(c *gin.Context) {
			ReverseProxy(c, "http://127.0.0.1:5173")
		})
	} else {
		r.GET("/", serveIndex)
		subFS, _ := fs.Sub(EmbeddedWebFS, "web/build")
		r.NoRoute(gin.WrapH(http.FileServer(http.FS(subFS))))
	}
}

func serveIndex(c *gin.Context) {
	file, err := EmbeddedWebFS.ReadFile("web/build/index.html")
	if err != nil {
		c.String(500, "Failed to read file")
		return
	}
	c.Data(200, "text/html", file)
}
