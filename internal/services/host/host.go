package host

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/nodeboxhq/nodebox-dashboard/internal/cmd"
	"github.com/nodeboxhq/nodebox-dashboard/internal/db/models"
	"github.com/nodeboxhq/nodebox-dashboard/internal/logger"
	"github.com/nodeboxhq/nodebox-dashboard/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

type UpdateInfo struct {
	Error       string `json:"error"`
	FileName    string `json:"fileName"`
	DownloadURL string `json:"downloadUrl"`
	Version     string `json:"version"`
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
	store.Version = cmd.Version

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

func (s *Service) AutoUpdateDashboard() {
	updateURL := "https://updater.nodebox.cloud/"
	resp, err := http.Get(updateURL)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var updateInfo UpdateInfo
	if err := json.Unmarshal(body, &updateInfo); err != nil {
		return
	}

	if updateInfo.Error != "" {
		return
	}

	if updateInfo.FileName == "" || updateInfo.DownloadURL == "" || updateInfo.Version == "" {
		return
	}

	currentVersion, err := semver.NewVersion(cmd.Version)
	if err != nil {
		return
	}

	newVersion, err := semver.NewVersion(updateInfo.Version)
	if err != nil {
		return
	}

	if currentVersion.Equal(newVersion) || currentVersion.GreaterThan(newVersion) {
		return
	}

	resp, err = http.Get(updateInfo.DownloadURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	newBinaryPath := filepath.Join(os.TempDir(), updateInfo.FileName)
	out, err := os.Create(newBinaryPath)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	currentBinaryPath := "/opt/nodebox-dashboard/nodebox-dashboard"
	if err := os.Rename(newBinaryPath, currentBinaryPath); err != nil {
		return
	}

	if err := os.Chmod(currentBinaryPath, 0755); err != nil {
		return
	}

	cmd := exec.Command("systemctl", "restart", "nodebox-dashboard")

	if err := cmd.Run(); err != nil {
		return
	}
}

func (s *Service) StartUpdateChecker(ctx context.Context) {
	if runtime.GOOS != "linux" {
		return
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.AutoUpdateDashboard()
		case <-ctx.Done():
			logger.L.Info().Msg("Stopping update checker")
			return
		}
	}
}

func (s *Service) InstallService() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("this function is only supported on Linux")
	}

	if _, err := exec.LookPath("systemctl"); err != nil {
		return fmt.Errorf("systemd is not installed on this system")
	}

	serviceName := "nodebox-dashboard.service"
	cmd := exec.Command("systemctl", "is-active", serviceName)
	if err := cmd.Run(); err == nil {
		return fmt.Errorf("service is already installed and active")
	}

	serviceContent := `[Unit]
Description=Nodebox Dashboard Service
After=network.target

[Service]
WorkingDirectory=/opt/nodebox-dashboard
ExecStart=/opt/nodebox-dashboard/nodebox-dashboard
Restart=always
User=root

[Install]
WantedBy=multi-user.target
`

	serviceFilePath := "/etc/systemd/system/" + serviceName
	err := os.WriteFile(serviceFilePath, []byte(serviceContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create service file: %v", err)
	}

	cmd = exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}

	cmd = exec.Command("systemctl", "enable", serviceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable service: %v, output: %s", err, string(output))
	}

	fmt.Println("Service installed successfully. Please restart the program using 'systemctl start nodebox-dashboard'")

	os.Exit(0)

	return nil
}
