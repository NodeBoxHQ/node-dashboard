package linea

import (
	"bytes"
	"encoding/json"
	"github.com/NodeboxHQ/node-dashboard/utils"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Status struct {
	Failure       bool `json:"failure"`
	Syncing       bool `json:"syncing"`
	CurrentHeight int  `json:"currentHeight"`
	MaxHeight     int  `json:"maxHeight"`
}

func NodeStatus() Status {
	status := Status{}
	customLineaNodeIP := utils.GetCustomLineaNodeIP()

	var url string

	if customLineaNodeIP != "" {
		url = "http://" + customLineaNodeIP + ":8545"
	} else {
		url = "http://127.0.0.1:8545"
	}

	syncStatusReq := []byte(`{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}`)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(syncStatusReq))
	if err != nil {
		status.Failure = true
		return status
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		status.Failure = true
		return status
	}

	var syncResult map[string]interface{}
	err = json.Unmarshal(body, &syncResult)
	if err != nil {
		return Status{}
	}

	if result, ok := syncResult["result"].(bool); ok && !result {
		status.Syncing = false
	} else {
		status.Syncing = true

		if result, ok := syncResult["result"].(map[string]interface{}); ok {
			if currentBlock, ok := result["currentBlock"].(string); ok {
				currentHeight, _ := strconv.ParseInt(currentBlock, 0, 64)
				status.CurrentHeight = int(currentHeight)
			}
			if highestBlock, ok := result["highestBlock"].(string); ok {
				maxHeight, _ := strconv.ParseInt(highestBlock, 0, 64)
				status.MaxHeight = int(maxHeight)
			}
		}
	}

	if !status.Syncing {
		blockNumberReq := []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(blockNumberReq))
		if err != nil {
			status.Failure = true
			return status
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			status.Failure = true
			return status
		}

		var blockNumberResult map[string]interface{}
		err = json.Unmarshal(body, &blockNumberResult)

		if err != nil {
			return Status{
				Failure: true,
			}
		}

		if result, ok := blockNumberResult["result"].(string); ok {
			currentHeight, _ := strconv.ParseInt(result, 0, 64)
			status.CurrentHeight = int(currentHeight)
		}
	}

	return status
}
