package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type RPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	ID      int         `json:"id"`
}

type NetworkNameResponse struct {
	NetworkName string `json:"networkName"`
}

type NodeInfoResponse struct {
	NumPeers string `json:"numPeers"`
	Peers    []struct {
		NodeID         string `json:"nodeID"`
		Version        string `json:"version"`
		LastSent       string `json:"lastSent"`
		LastReceived   string `json:"lastReceived"`
		ObservedUptime string `json:"observedUptime"`
	} `json:"peers"`
}

type Juneo struct {
	NodeID           string  `json:"nodeId"`
	Status           string  `json:"status"`
	UptimePercentage float64 `json:"uptimePercentage"`
	NetworkName      string  `json:"networkName"`
}

func getNodeID() (string, error) {
	response, err := http.Post("http://127.0.0.1:9650/ext/info", "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"info.getNodeID"}`))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var result struct {
		Result struct {
			NodeID string `json:"nodeID"`
		} `json:"result"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Result.NodeID, nil
}

func getSyncStatus() bool {
	response, err := http.Post("http://127.0.0.1:9650/ext/info", "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"info.isBootstrapped","params":{"chain":"JUNE"}}`))

	if err != nil {
		return false
	}

	defer response.Body.Close()

	var result map[string]interface{}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return false
	}

	if isBootstrapped, ok := result["result"].(map[string]interface{})["isBootstrapped"].(bool); ok {
		return isBootstrapped
	}

	return false
}

func getNodeInfo(nodeID string) (NodeInfoResponse, error) {
	payload := fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"id":1,
		"method":"info.peers",
		"params": {
			"nodeIDs": ["%s"]
		}
	}`, nodeID)

	response, err := http.Post("https://rpc.juneo-mainnet.network/ext/info", "application/json", strings.NewReader(payload))
	if err != nil {
		return NodeInfoResponse{}, err
	}
	defer response.Body.Close()

	var rpcResponse RPCResponse
	if err := json.NewDecoder(response.Body).Decode(&rpcResponse); err != nil {
		return NodeInfoResponse{}, err
	}

	var nodeInfo NodeInfoResponse
	resultJSON, err := json.Marshal(rpcResponse.Result)
	if err != nil {
		return NodeInfoResponse{}, err
	}

	if err := json.Unmarshal(resultJSON, &nodeInfo); err != nil {
		return NodeInfoResponse{}, err
	}

	return nodeInfo, nil
}

func JuneoInfo() Juneo {
	juneo := Juneo{
		NodeID:           "Unknown",
		Status:           "Offline",
		UptimePercentage: 0.0,
		NetworkName:      "Unknown",
	}

	if getSyncStatus() {
		juneo.Status = "Online"

		response, err := http.Post("http://127.0.0.1:9650/ext/info", "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"info.getNetworkName"}`))
		if err == nil {
			defer response.Body.Close()
			if response.StatusCode == http.StatusOK {
				var result struct {
					JSONRPC string              `json:"jsonrpc"`
					Result  NetworkNameResponse `json:"result"`
					ID      int                 `json:"id"`
				}
				if err := json.NewDecoder(response.Body).Decode(&result); err == nil {
					juneo.NetworkName = result.Result.NetworkName
				}
			}
		}

		nodeID, err := getNodeID()
		if err == nil {
			juneo.NodeID = nodeID

			nodeInfo, err := getNodeInfo(nodeID)
			if err == nil {
				if len(nodeInfo.Peers) > 0 {
					uptimeFloat := 0.0
					fmt.Sscanf(nodeInfo.Peers[0].ObservedUptime, "%f", &uptimeFloat)
					juneo.UptimePercentage = uptimeFloat
				}
			}
		}
	}

	return juneo
}
