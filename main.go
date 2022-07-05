package main

import (
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	"log"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	configs := readConfig()
	items := createItems(configs)

	index, err := fuzzyfinder.Find(
		items,
		itemFunc(items),
		fuzzyfinder.WithPreviewWindow(previewFunc(items)))
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
	resultChannel := make(chan []Item, len(configs))
	var wg sync.WaitGroup
	for _, config := range configs {
		wg.Add(1)
		go parseItemsFromConfig(config, &wg, resultChannel)
	}
	wg.Wait()
	for i := 0; i < len(configs); i++ {
		items = append(items, <-resultChannel...)
	}
	return items
}

func parseItemsFromConfig(config ConfigItem, wg *sync.WaitGroup, channel chan []Item) {
	defer wg.Done()
	listCommand := exec.Command("bash", "-c", config.List)
	output, err2 := listCommand.Output()
	if err2 != nil {
		log.Fatal(err2)
	}
	var items []Item
	for _, line := range strings.Split(strings.TrimSuffix(string(output), "\n"), "\n") {
		parts := strings.Split(line, ",")
		item := Item{
			Name:    parts[0],
			Command: config.Command,
			Path:    parts[1],
		}
		items = append(items, item)
	}
	channel <- items
}

type Item struct {
	Name    string
	Command string
	Path    string
}

func (item *Item) DisplayName() string {
	return item.Command + " " + item.Name
}

func itemFunc(items []Item) func(index int) string {
	return func(index int) string {
		return items[index].DisplayName()
	}
}

func previewFunc(items []Item) func(i int, w int, h int) string {
	return func(i, w, h int) string {
		if i == -1 {
			return ""
		}
		return fmt.Sprintf("opening: %s", items[i].Path)
	}
}
