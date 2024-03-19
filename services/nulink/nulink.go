package nulink

import "github.com/NodeboxHQ/node-dashboard/utils"

type Status struct {
	Online bool `json:"online"`
}

func NodeStatus() Status {
	status := Status{
		Online: true,
	}

	inUse := utils.IsPortInUse(9151)

	if !inUse {
		status.Online = false
		return status
	}

	return status
}
