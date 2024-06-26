package xally

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/NodeboxHQ/node-dashboard/services/config"
	"github.com/NodeboxHQ/node-dashboard/utils"
	"github.com/NodeboxHQ/node-dashboard/utils/logger"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

const baseURL = "https://api-node.xally.ai"
const levelPath = "/root/.config/xally_client/Local Storage/leveldb/000003.log"

var (
	lock         sync.Mutex
	nodeData     []NodeInfo
	previousData []NodeInfo
)

type ApiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type NodeInfo struct {
	Chain       string  `json:"chain"`
	KeyID       string  `json:"key_id"`
	NodeID      string  `json:"node_id"`
	RunningTime float64 `json:"running_time"`
	Point       float64 `json:"point"`
	Status      string  `json:"status"`
	LastCheckTS int64   `json:"last_check_ts"`
}

func FetchNodeData() ([]NodeInfo, error) {
	lock.Lock()
	defer lock.Unlock()

	if len(nodeData) > 0 {
		if time.Now().Unix()-nodeData[0].LastCheckTS < 300 {
			return nodeData, nil
		}
	}

	jwtToken, err := getJwtToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", baseURL+"/nodes/info", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+jwtToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		jwtToken, err = getJwtToken()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+jwtToken)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var nodes []NodeInfo
	if err := json.Unmarshal(apiResp.Data, &nodes); err != nil {
		return nil, err
	}

	lastCheckTs := time.Now().Unix()
	for i := range nodes {
		nodes[i].LastCheckTS = lastCheckTs
	}

	previousData = nodeData
	nodeData = nodes

	return nodes, nil
}

func CheckRunning(config *config.Config) {
	nodes, err := FetchNodeData()
	allGood := true
	if err != nil {
		logger.Error("Error fetching node data: ", err)
		allGood = false
	}

	for _, node := range nodes {
		if node.Status == "running" {
			continue
		} else {
			logger.Error("Node is not running: ", node.NodeID)
			allGood = false
			break
		}
	}

	previousDataMap := make(map[string]NodeInfo)
	for _, node := range previousData {
		previousDataMap[node.NodeID] = node
	}

	var unchangedNodes []string

	for _, node := range nodes {
		if prevNode, ok := previousDataMap[node.NodeID]; ok {
			if node.RunningTime == prevNode.RunningTime && node.Point == prevNode.Point {
				unchangedNodes = append(unchangedNodes, node.NodeID)
			}
		}
	}

	if len(unchangedNodes) > 0 {
		logger.Error("Nodes have not changed: ", unchangedNodes)
		allGood = false
	}

	if !allGood {
		timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
		message := "-----------------------------------------------------------------------\n"
		message += fmt.Sprintf("🚨 **Xally** Node Malfunction Detected @ %s UTC 🚨\n\n", timestamp)
		message += fmt.Sprintf("👑 **Owner**: %s\n", config.Owner)
		message += fmt.Sprintf("🌐 **IPv4**: %s\n", config.IPv4)
		message += fmt.Sprintf("🌐 **IPv6**: %s\n", config.IPv6)
		message += fmt.Sprintf("⚙️ **Dashboard Version**: %s\n", config.NodeboxDashboardVersion)
		message += "-----------------------------------------------------------------------\n"

		utils.SendAlert(message)
	} else {
		logger.Info("All Xally nodes are running")
	}
}

func getJwtToken() (string, error) {
	file, err := os.Open(levelPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var lastToken string
	tokenRegex := regexp.MustCompile(`[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+`)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := tokenRegex.FindAllString(line, -1)
		if len(tokens) > 0 {
			lastToken = tokens[len(tokens)-1] // Get the last token in the line
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan the file: %v", err)
	}

	if lastToken == "" {
		return "", fmt.Errorf("no JWT token found in the file")
	}

	return lastToken, nil
}
