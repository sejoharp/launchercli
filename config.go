package main

import (
	"encoding/json"
	"log"
	"os"
)

func readAndParseConfigFile(pathToConfigFile string) []ConfigItem {
	file, err := os.ReadFile(pathToConfigFile)
	if err != nil {
		log.Fatalf("unable to read config file, error: %v", err)
	}
	var configItems []ConfigItem
	if err := json.Unmarshal(file, &configItems); err != nil {
		log.Fatalf("unable to parse config file, error: %v", err)
	}
	return configItems
}

type ConfigItem struct {
	DisplayName string `json:"display_name"`
	List        string `json:"list"`
	Command     string `json:"command"`
}
