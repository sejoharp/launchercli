package main

import (
	"encoding/json"
	"log"
	"os"
)

func readConfig() []ConfigItem {
	file, err := os.ReadFile(homeDir() + "/.launcher-config.json")
	if err != nil {
		log.Fatalf("unable to read config file, error: %v", err)
	}
	var configItems []ConfigItem
	if err := json.Unmarshal(file, &configItems); err != nil {
		log.Fatalf("unable to parse config file, error: %v", err)
	}
	return configItems
}

func homeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return dirname
}

type ConfigItem struct {
	List    string `json:"list"`
	Command string `json:"command"`
}
