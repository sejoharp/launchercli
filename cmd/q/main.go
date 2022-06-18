package main

import (
	"encoding/json"
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func homeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return dirname
}

func main() {
	configs := readConfig()
	var items []Item
	for _, config := range configs {
		listCommand := exec.Command("bash", "-c", config.List)
		output, err2 := listCommand.Output()
		if err2 != nil {
			log.Fatal(err2)
		}
		for _, line := range strings.Split(strings.TrimSuffix(string(output), "\n"), "\n") {
			item := Item{
				Name:    filepath.Base(line),
				Command: config.Command,
				Path:    line,
			}
			items = append(items, item)
		}
	}

	//fmt.Println(items)
	//fmt.Println("ende")
	//os.Exit(999)
	//var repos []Item
	//repos = append(repos, listDirs(homeDir()+"/private/", "code")...)
	//repos = append(repos, listDirs(homeDir()+"/repos/", "code")...)
	//repos = append(repos, listDirs(homeDir()+"/private/", "idea")...)
	//repos = append(repos, listDirs(homeDir()+"/repos/", "idea")...)

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

type ConfigItem struct {
	List                  string `json:"list"`
	LastPartAsDisplayName bool   `json:"lastPartAsDisplayName"`
	Command               string `json:"command"`
}
type Item struct {
	Name    string
	Command string
	Path    string
}

func (item *Item) DisplayName() string {
	return item.Command + " " + item.Name
}

func readConfig() []ConfigItem {
	file, err := os.ReadFile(homeDir() + "/.launcher-config.json")
	if err != nil {
		log.Fatalf("unable to read config file, error: %v", err)
	}
	var configItems []ConfigItem
	if err := json.Unmarshal(file, &configItems); err != nil {
		panic(err)
	}
	return configItems
}

func listDirs(path string, command string) []Item {
	dirHandle, err := os.Open(path)
	if err != nil {
		log.Fatalf("directory not found %s, error %v", path, err)
		return nil
	}
	dirs, err := dirHandle.Readdirnames(-1)
	if err != nil {
		log.Fatalf("error reading directory for %s, error %v", path, err)
		return []Item{}
	}
	var items []Item
	for _, dir := range dirs {
		if strings.HasPrefix(dir, ".") {
			continue
		}
		items = append(items, Item{
			Name:    dir,
			Path:    filepath.Join(path, dir),
			Command: command,
		})
	}
	return items
}
