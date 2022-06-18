package main

import (
	"encoding/json"
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	configs := readConfig()
	items := createItems(configs)

	index, err := fuzzyfinder.Find(
		items,
		itemFunc(items),
		fuzzyfinder.WithPreviewWindow(preview(items)))
	if err != nil {
		log.Fatal(err)
	}
	cliCommand := exec.Command(items[index].Command, items[index].Path)
	err = cliCommand.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func createItems(configs []ConfigItem) []Item {
	var items []Item
	for _, config := range configs {
		listCommand := exec.Command("bash", "-c", config.List)
		output, err2 := listCommand.Output()
		if err2 != nil {
			log.Fatal(err2)
		}
		for _, line := range strings.Split(strings.TrimSuffix(string(output), "\n"), "\n") {
			parts := strings.Split(line, ",")
			item := Item{
				Name:    parts[0],
				Command: config.Command,
				Path:    parts[1],
			}
			items = append(items, item)
		}
	}
	return items
}

func itemFunc(items []Item) func(i int) string {
	return func(i int) string {
		return items[i].DisplayName()
	}
}

func preview(items []Item) func(i int, w int, h int) string {
	return func(i, w, h int) string {
		if i == -1 {
			return ""
		}
		return fmt.Sprintf("opening: %s", items[i].Path)
	}
}

type Item struct {
	Name    string
	Command string
	Path    string
}

func (item *Item) DisplayName() string {
	return item.Command + " " + item.Name
}

type ConfigItem struct {
	List    string `json:"list"`
	Command string `json:"command"`
}

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
