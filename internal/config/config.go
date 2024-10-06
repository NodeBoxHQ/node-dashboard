package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nodeboxhq/nodebox-dashboard/internal"
)

var ParsedConfig *internal.NodeboxConfig

func ParseConfig(path string) *internal.NodeboxConfig {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	ParsedConfig = &internal.NodeboxConfig{}
	err = decoder.Decode(ParsedConfig)

	if err != nil {
		log.Fatal(err)
	}

	err = SetupDataPath(ParsedConfig)

	if err != nil {
		log.Fatal(err)
	}

	return ParsedConfig
}

func SetupDataPath(cfg *internal.NodeboxConfig) error {
	if cfg.DataPath == "" {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			return err
		}

		cfg.DataPath = homeDir + "/.nodebox"

		if _, err := os.Stat(cfg.DataPath); os.IsNotExist(err) {
			err := os.Mkdir(cfg.DataPath, 0755)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetLineaIP() string {
	if ParsedConfig.LineaIP == "" {
		return "127.0.0.1"
	}

	return ParsedConfig.LineaIP
}
