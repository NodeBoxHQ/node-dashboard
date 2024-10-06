package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/nodeboxhq/nodebox-dashboard/internal/config"
)

type Linea struct {
	Status        string `json:"status"`
	CurrentHeight int    `json:"currentHeight"`
	MaxHeight     int    `json:"maxHeight"`
}

func LineaInfo() Linea {
	var linea Linea

	nodeIp := config.GetLineaIP()
	url := fmt.Sprintf("http://%s:8545", nodeIp)

	syncStatusReq := []byte(`{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}`)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(syncStatusReq))
	if err != nil {
		linea.Status = "Offline"
		return linea
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		linea.Status = "Offline"
		return linea
	}

	var syncResult map[string]interface{}
	err = json.Unmarshal(body, &syncResult)
	if err != nil {
		linea.Status = "Offline"
		return linea
	}

	if result, ok := syncResult["result"].(bool); ok && !result {
		linea.Status = "Online"
	} else {
		if result, ok := syncResult["result"].(map[string]interface{}); ok {
			if currentBlock, ok := result["currentBlock"].(string); ok {
				currentHeight, _ := strconv.ParseInt(currentBlock, 0, 64)
				linea.CurrentHeight = int(currentHeight)
			}

			if highestBlock, ok := result["highestBlock"].(string); ok {
				maxHeight, _ := strconv.ParseInt(highestBlock, 0, 64)
				linea.MaxHeight = int(maxHeight)
			}
		}

		linea.Status = "Syncing"
	}

	if linea.Status == "Online" {
		blockNumberReq := []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(blockNumberReq))

		if err != nil {
			linea.Status = "Offline"
			return linea
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			linea.Status = "Offline"
			return linea
		}

		var blockNumberResult map[string]interface{}
		err = json.Unmarshal(body, &blockNumberResult)

		if err != nil {
			linea.Status = "Offline"
			return linea
		}

		if result, ok := blockNumberResult["result"].(string); ok {
			currentHeight, _ := strconv.ParseInt(result, 0, 64)
			linea.CurrentHeight = int(currentHeight)
			linea.MaxHeight = int(currentHeight)
			linea.Status = "Online"
		}
	}

	return linea
}
