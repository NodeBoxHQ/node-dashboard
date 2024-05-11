package xally

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

const baseURL = "https://api-node.xally.ai"
const levelPath = "/root/.config/xally_client/Local Storage/leveldb/000003.log"

var (
	lock sync.Mutex
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

var nodeData []NodeInfo

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

	nodeData = nodes

	return nodes, nil
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
