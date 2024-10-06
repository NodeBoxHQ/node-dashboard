package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func ReverseProxy(c *gin.Context, backend string) {
	remote, err := url.Parse(backend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse proxy URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		if !strings.Contains(err.Error(), "context canceled") {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		}
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
