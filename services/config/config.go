package config

import (
	"encoding/json"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"os"
)

const version = "1.0.0"

type Config struct {
	Node  string `json:"node"`
	IPv4  string `json:"ip"`
	IPv6  string `json:"ip6"`
	Owner string `json:"owner"`
	Port  int    `json:"port"`
}

var loadedConfig *Config

func LoadConfig() (*Config, error) {
	path := "./config.json"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)

	if err != nil {
		return nil, err
	}

	loadedConfig = &config

	return &config, nil
}

func ShowAsciiArt() {
	myFigure := figure.NewFigure("Nodebox", "doom", true)
	myFigure.Print()
	fmt.Println("\n\t\t\t\t\tVersion: ", version)
	fmt.Println("\n")
}
