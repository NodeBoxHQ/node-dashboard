package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Hyperliquid struct {
	Version       string    `json:"version"`
	Status        string    `json:"status"`
	Height        int64     `json:"currentHeight"`
	BlockTime     time.Time `json:"blockTime"`
	ApplyDuration float64   `json:"applyDuration"`
}

func HyperliquidInfo() (Hyperliquid, error) {
	var hl Hyperliquid

	hl.Status = "Offline"

	cmd := exec.Command("/root/hl-node", "--version")
	output, err := cmd.Output()
	if err != nil {
		return hl, fmt.Errorf("failed to get version: %v", err)
	}

	hl.Version = strings.TrimSpace(string(output))

	blockTimesDir := "/root/hl/data/block_times"
	files, err := os.ReadDir(blockTimesDir)
	if err != nil {
		return hl, fmt.Errorf("failed to read block_times directory: %v", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	if len(files) > 0 {
		oldestFile := filepath.Join(blockTimesDir, files[0].Name())

		file, err := os.Open(oldestFile)
		if err != nil {
			return hl, fmt.Errorf("failed to open file %s: %v", oldestFile, err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("failed to close file %s: %v", oldestFile, err)
			}
		}(file)

		var lastLine string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lastLine = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return hl, fmt.Errorf("error reading file %s: %v", oldestFile, err)
		}

		var blockData struct {
			Height        int64     `json:"height"`
			BlockTime     time.Time `json:"block_time"`
			ApplyDuration float64   `json:"apply_duration"`
		}

		if err := json.Unmarshal([]byte(lastLine), &blockData); err != nil {
			return hl, fmt.Errorf("failed to parse JSON: %v", err)
		}

		hl.Status = "Online"
		hl.Height = blockData.Height
		hl.BlockTime = blockData.BlockTime
		hl.ApplyDuration = blockData.ApplyDuration
	}

	return hl, nil
}
