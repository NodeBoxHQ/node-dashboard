package node

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

type Dusk struct {
	Status  string
	Height  int
	Version string
}

type InfoResponse struct {
	Version string `json:"version"`
}

func DuskInfo() Dusk {
	dusk := Dusk{
		Status:  "Offline",
		Height:  0,
		Version: "Unknown",
	}

	cmd := exec.Command("ruskquery", "block-height")
	output, err := cmd.Output()
	if err != nil {
		return dusk
	}

	outputStr := strings.TrimSpace(string(output))
	height, err := strconv.Atoi(outputStr)
	if err != nil {
		return dusk
	}

	dusk.Status = "Online"
	dusk.Height = height

	cmd = exec.Command("ruskquery", "info")
	output, err = cmd.Output()
	if err != nil {
		return dusk
	}

	var info InfoResponse
	err = json.Unmarshal(output, &info)
	if err != nil {
		dusk.Version = "Unknown"
	} else {
		dusk.Version = info.Version
	}

	return dusk
}
