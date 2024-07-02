package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/services/config"
	"github.com/NodeboxHQ/node-dashboard/utils/logger"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func SecondsToReadable(seconds int) string {
	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	timeString := ""

	if days > 0 {
		timeString += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 || days > 0 { // Include hours if there are any days
		timeString += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 || hours > 0 || days > 0 { // Include minutes if there are any hours or days
		timeString += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 || minutes > 0 || hours > 0 || days > 0 { // Include seconds if there are any minutes, hours, or days
		timeString += fmt.Sprintf("%ds", seconds)
	}

	return timeString
}

func InstallService() {
	hostname, _ := os.Hostname()

	if hostname == "hayzam-pc" || hostname == "nodebox-xally" {
		return
	}

	if strings.Contains(hostname, "-nulink") {
		return
	}

	if _, err := os.Stat("/etc/systemd/system/nodebox-dashboard.service"); os.IsNotExist(err) {
		var service string

		service = `
[Unit]
Description=Nodebox Dashboard
After=network.target

[Service]
User=root
Group=root
Type=simple
Restart=always
RestartSec=5
WorkingDirectory=/opt/nodebox-dashboard
ExecStart=/opt/nodebox-dashboard/nodebox-dashboard

[Install]
WantedBy=multi-user.target
`

		err := os.WriteFile("/etc/systemd/system/nodebox-dashboard.service", []byte(service), 0644)

		if err != nil {
			logger.Error("Error writing service file", err)
		}

		_, err = exec.Command("systemctl", "daemon-reload").Output()

		if err != nil {
			logger.Error("Error reloading systemd daemon", err)
		}

		_, err = exec.Command("systemctl", "enable", "nodebox-dashboard").Output()

		if err != nil {
			logger.Error("Error enabling service", err)
		}

		logger.Info("Service installed, exiting now...run 'systemctl start nodebox-dashboard' to start the service")

		os.Exit(0)
	} else {
		logger.Info("Service already installed")
	}
}

func IsPortInUse(port int) bool {
	tcpAddr, tcpErr := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if tcpErr != nil {
		fmt.Println("Error resolving TCP address:", tcpErr)
		return true
	}
	tcpLn, tcpErr := net.ListenTCP("tcp", tcpAddr)
	if tcpErr != nil {
		return true
	}
	defer tcpLn.Close()

	udpAddr, udpErr := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if udpErr != nil {
		fmt.Println("Error resolving UDP address:", udpErr)
		return true
	}
	udpLn, udpErr := net.ListenUDP("udp", udpAddr)
	if udpErr != nil {
		return true
	}
	defer udpLn.Close()

	return false
}

func GetCurrentCPUArch() string {
	cmd := exec.Command("uname", "-m")
	stdout, err := cmd.Output()

	if err != nil {
		logger.Error("Error getting CPU architecture", err)
		return ""
	}

	return strings.TrimSpace(string(stdout))
}

func DownloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

func SelfUpdate(version string) {
	nodeType := config.GetNodeType()

	if nodeType == "nulink" {
		return
	}

	version = strings.Replace(version, ".", "", -1)
	versionInt, err := strconv.Atoi(version)

	if err != nil {
		logger.Error("Error converting version to integer:", err)
		return
	}

	url := FetchLatestReleaseDownloadURL("NodeBoxHQ", "nodebox-dashboard", versionInt)

	if url == "" {
		return
	}

	logger.Info("New version found, downloading from:", url)

	err = DownloadFile(url, "/tmp/nodebox-dashboard")

	if err != nil {
		logger.Error("Error downloading new version:", err)
		return
	}

	newBinary := "/tmp/nodebox-dashboard"
	oldBinary := "/opt/nodebox-dashboard/nodebox-dashboard"

	err = os.Rename(newBinary, oldBinary)

	if err != nil {
		logger.Error("Error replacing old binary with new binary:", err)
		return
	}

	_, err = exec.Command("chmod", "+x", oldBinary).Output()

	if err != nil {
		logger.Error("Error changing permissions of new binary:", err)
		return
	}

	_, err = exec.Command("systemctl", "restart", "nodebox-dashboard").Output()

	if err != nil {
		logger.Error("Error restarting service:", err)
		return
	}
}

type Release struct {
	TagName            string `json:"tag_name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Assets             []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func FetchLatestReleaseDownloadURL(owner, repo string, version int) string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error fetching latest release:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Error fetching latest release:", resp.Status)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return ""
	}

	var releases []Release
	if err := json.Unmarshal(body, &releases); err != nil {
		logger.Error("Error unmarshalling response body:", err)
		return ""
	}

	if len(releases) > 0 {
		latestRelease := releases[0]
		for _, asset := range latestRelease.Assets {
			if strings.Contains(asset.Name, GetCurrentCPUArch()) {
				pattern := `\d+\.\d+\.\d+`
				re := regexp.MustCompile(pattern)
				maybeLatest := re.FindString(asset.Name)
				maybeLatestInt, err := strconv.Atoi(strings.Replace(maybeLatest, ".", "", -1))

				if err != nil {
					logger.Error("Error converting maybeLatest to integer:", err)
					return ""
				}

				if maybeLatestInt > version {
					return asset.BrowserDownloadURL
				}
			}
		}
	} else {
		logger.Error("No releases found")
		return ""
	}

	logger.Debug("No new version found")
	return ""
}

type WebhookMessage struct {
	Content string `json:"content"`
}

func SendAlert(content string) bool {
	webhookURL := "https://discord.com/api/webhooks/1239543600816324669/HtvV0owSKPXAFiMpbAAVnEu2Q28kh84nwyl-uIbAFr4N8nYtj0Nd8MeQcf5036hgIbBu"
	message := WebhookMessage{Content: content}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		logger.Error("Error marshalling alert message:", err)
		return false
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		logger.Error("Error sending alert message:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		logger.Error("Error sending alert message:", resp.Status)
		return false
	}

	return true
}

func GetIPForAccess(ip string) string {
	if _, err := os.Stat("/root/.nodebox-ip"); os.IsNotExist(err) {
		return ip
	} else {
		rip, err := os.ReadFile("/root/.nodebox-ip")
		if err != nil {
			logger.Error("Error reading .nodebox-ip file:", err)
			return ip
		}
		return strings.TrimSpace(string(rip))
	}
}
