package xally

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

const baseURL = "https://api-node.xally.ai"

var (
	apiKey     string
	authToken  string
	retryCount int
	lock       sync.Mutex
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

func GetXallyAPIKey() string {
	file, err := os.Open("/root/.config/xally_client/Local Storage/leveldb/000003.log")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	regex := regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := regex.FindString(line); matches != "" {
			return matches
		}
	}
	return ""
}

func FetchNodeData() ([]NodeInfo, error) {
	lock.Lock()
	defer lock.Unlock()

	if len(nodeData) > 0 {
		if time.Now().Unix()-nodeData[0].LastCheckTS < 300 {
			return nodeData, nil
		}
	}

	var nodes []NodeInfo
	backoff := 1 * time.Second
	const maxBackoff = 10 * time.Minute
	retryCount = 0

	for {
		if retryCount > 5 {
			fmt.Println("Max retries exceeded")
			return nil, fmt.Errorf("max retries exceeded")
		}

		req, err := http.NewRequest("GET", baseURL+"/nodes/info", nil)
		if err != nil {
			return nil, err
		}

		req.Header.Add("Authorization", "Bearer "+authToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			time.Sleep(backoff)
			backoff = min(2*backoff, maxBackoff)
			retryCount++
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 401 {
			newToken, err := getAuthKey(apiKey)
			if err != nil {
				return nil, err
			}
			authToken = newToken
			continue
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
		}

		var apiResp ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			time.Sleep(backoff)
			backoff = min(2*backoff, maxBackoff)
			retryCount++
			continue
		}

		if err := json.Unmarshal(apiResp.Data, &nodes); err != nil {
			return nil, err
		}

		lastCheckTs := time.Now().Unix()

		for i := range nodes {
			nodes[i].LastCheckTS = lastCheckTs
		}

		nodeData = nodes
		break
	}

	return nodes, nil
}

func getAuthKey(apiKey string) (string, error) {
	if apiKey == "" {
		apiKey = GetXallyAPIKey()
	}

	payload := fmt.Sprintf(`{"api_key":"%s"}`, apiKey)
	req, err := http.NewRequest("POST", baseURL+"/auth/api-key", bytes.NewBufferString(payload))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}

	if apiResp.Code != 2000 {
		return "", fmt.Errorf("failed to refresh auth key: %s", apiResp.Message)
	}

	key := make(map[string]interface{})

	if err := json.Unmarshal(apiResp.Data, &key); err != nil {
		return "", err
	}

	return key["access_token"].(string), nil
}
