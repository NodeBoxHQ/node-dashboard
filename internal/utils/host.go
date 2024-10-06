package utils

import (
	"os/exec"
	"strconv"
	"strings"
)

func GetUptimeSeconds() (int64, error) {
	out, err := exec.Command("cat", "/proc/uptime").Output()
	if err != nil {
		return 0, err
	}
	uptimeStr := strings.Fields(string(out))[0]
	uptimeSeconds, err := strconv.ParseFloat(uptimeStr, 64)
	if err != nil {
		return 0, err
	}
	return int64(uptimeSeconds), nil
}
