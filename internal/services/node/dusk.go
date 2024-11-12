package node

import (
	"encoding/json"
	"github.com/nodeboxhq/nodebox-dashboard/internal/config"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

type Dusk struct {
	Status  string `json:"status"`
	Height  int    `json:"currentHeight"`
	Version string `json:"version"`
	Stake   Stake  `json:"stake"`
}

type Info struct {
	Version string `json:"version"`
}

type Stake struct {
	StakingAddress string  `json:"stakingAddress"`
	EligibleStake  float64 `json:"eligibleStake"`
	Slashes        int     `json:"slashes"`
	HardSlashes    int     `json:"hardSlashes"`
	Rewards        float64 `json:"rewards"`
}

func GetStakeInfo() Stake {
	cmd := exec.Command("rusk-wallet", "--password", config.GetDuskPassword(), "stake-info")
	output, err := cmd.Output()
	if err != nil {
		return Stake{}
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	stake := Stake{}

	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		value := strings.TrimSpace(parts[1])

		switch {
		case strings.Contains(line, "Staking address:"):
			stake.StakingAddress = value
		case strings.Contains(line, "Eligible stake:"):
			value = strings.Split(value, " ")[0]
			if s, err := strconv.ParseFloat(value, 64); err == nil {
				stake.EligibleStake = math.Round(s*100) / 100
			}
		case strings.Contains(line, "Reclaimable slashed stake:"):
			value = strings.Split(value, " ")[0]
			stake.Slashes, _ = strconv.Atoi(value)
		case strings.Contains(line, "Hard Slashes:"):
			stake.HardSlashes, _ = strconv.Atoi(value)
		case strings.Contains(line, "Accumulated rewards is:"):
			value = strings.Split(value, " ")[0]
			if s, err := strconv.ParseFloat(value, 64); err == nil {
				stake.Rewards = math.Round(s*100) / 100
			}
		}
	}

	return stake
}

func DuskInfo() Dusk {
	dusk := Dusk{
		Status:  "Offline",
		Height:  0,
		Version: "Unknown",
		Stake:   GetStakeInfo(),
	}

	cmd := exec.Command("ruskquery", "block-height")
	output, _ := cmd.Output()
	outputStr := strings.TrimSpace(string(output))
	height, err := strconv.Atoi(outputStr)
	if err != nil {
		dusk.Height = 0
	}

	cmd = exec.Command("ruskquery", "info")
	output, err = cmd.Output()
	if err != nil {
		return dusk
	}

	var info Info
	err = json.Unmarshal(output, &info)
	if err != nil {
		dusk.Version = "Unknown"
	} else {
		dusk.Version = info.Version
	}

	if dusk.Version != "Unknown" {
		dusk.Status = "Online"
	}

	dusk.Height = height

	return dusk
}
