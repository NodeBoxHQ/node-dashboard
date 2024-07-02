package juneo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Status struct {
	NodeID            string  `json:"node_id"`
	Synced            bool    `json:"synced"`
	UptimePercentage  float64 `json:"uptime_percentage"`
	DelegationEndDate string  `json:"delegation_end_date"`
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

func getNodeID() string {
	response, err := http.Post("http://127.0.0.1:9650/ext/info", "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"info.getNodeID"}`))

	if err != nil {
		return "Unknown"
	}

	defer response.Body.Close()

	var result map[string]interface{}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "Unknown"
	}

	if resultNodeID, ok := result["result"].(map[string]interface{})["nodeID"].(string); ok {
		return resultNodeID
	}

	return "Unknown"
}

func getUptimePercentage() float64 {
	response, err := http.Post("http://127.0.0.1:9650/ext/info", "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"info.uptime"}`))
	if err != nil {
		return 0.0
	}
	defer response.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return 0.0
	}

	if _, ok := result["error"]; ok {
		return 0.0
	}

	resultData, ok := result["result"].(map[string]interface{})
	if !ok {
		return 0.0
	}

	rewardingStakePercentageStr, ok1 := resultData["rewardingStakePercentage"].(string)
	weightedAveragePercentageStr, ok2 := resultData["weightedAveragePercentage"].(string)
	if ok1 && ok2 {
		rewardingStakePercentageFloat, err1 := strconv.ParseFloat(rewardingStakePercentageStr, 64)
		weightedAveragePercentageFloat, err2 := strconv.ParseFloat(weightedAveragePercentageStr, 64)
		if err1 == nil && err2 == nil {
			return float64(int((rewardingStakePercentageFloat+weightedAveragePercentageFloat)/2*100)) / 100
		}
	}

	return 0.0
}

func getDelegationEndDate(nodeID string) string {
	response, err := http.Post("http://127.0.0.1:9650/ext/bc/P", "application/json", strings.NewReader(fmt.Sprintf(`{"jsonrpc":"2.0","method":"platform.getCurrentValidators","params":{"nodeIDs":["%s"]},"id":1}`, nodeID)))

	if err != nil {
		return "Unknown"
	}

	defer response.Body.Close()

	var result map[string]interface{}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "Unknown"
	}

	if validators, ok := result["result"].(map[string]interface{})["validators"].([]interface{}); ok {
		if len(validators) > 0 {
			if endTime, ok := validators[0].(map[string]interface{})["endTime"].(string); ok {
				endTimeInt, _ := strconv.ParseInt(endTime, 10, 64)
				endTimeUnix := time.Unix(endTimeInt, 0)
				return endTimeUnix.Format("02-01-2006")
			}
		}
	}

	return "Unknown"
}

func NodeStatus() Status {
	NodeID := getNodeID()

	status := Status{
		NodeID:            NodeID,
		Synced:            getSyncStatus(),
		UptimePercentage:  getUptimePercentage(),
		DelegationEndDate: getDelegationEndDate(NodeID),
	}

	return status
}
