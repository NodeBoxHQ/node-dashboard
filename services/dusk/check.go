package dusk

import (
	"os/exec"
	"strconv"
	"strings"
)

type Status struct {
	Failure bool `json:"failure"`
	Height  int  `json:"height"`
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

	return status
}
