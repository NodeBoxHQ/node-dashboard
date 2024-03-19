package dusk

import (
	"os/exec"
	"strconv"
	"strings"
)

type Status struct {
	Failure bool   `json:"failure"`
	Height  int    `json:"height"`
	Version string `json:"version"`
}

func NodeStatus() Status {
	cmd := exec.Command("ruskquery", "block-height")
	output, err := cmd.Output()

	status := Status{
		Failure: true,
		Height:  0,
	}

	if err != nil {
		return status
	}

	outputStr := strings.TrimSpace(string(output))
	height, err := strconv.Atoi(outputStr)

	if err != nil {
		return status
	}

	status.Height = height
	status.Failure = false

	cmd = exec.Command("ruskquery", "info")

	output, err = cmd.Output()

	if err != nil {
		status.Version = "Unknown"
	}

	if strings.Contains(string(output), "version") == false {
		status.Version = "Unknown"
	} else {
		status.Version = strings.Split(strings.Split(string(output), "version\": \"")[1], "\"")[0]
	}

	return status
}
