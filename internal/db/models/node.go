package models

type Dusk struct {
	ID            int    `json:"id"`
	Status        string `json:"status"`
	Version       string `json:"version"`
	CurrentHeight int    `json:"currentHeight"`
}
