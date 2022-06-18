package main

import (
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
	var repos []Item
	repos = append(repos, listDirs(homeDir()+"/private/", "code")...)
	repos = append(repos, listDirs(homeDir()+"/repos/", "code")...)
	repos = append(repos, listDirs(homeDir()+"/private/", "idea")...)
	repos = append(repos, listDirs(homeDir()+"/repos/", "idea")...)

	index, err := fuzzyfinder.Find(
		repos,
		itemFunc(repos),
		fuzzyfinder.WithPreviewWindow(preview(repos)))
	if err != nil {
		log.Fatal(err)
	}
	cliCommand := exec.Command(repos[index].Command, repos[index].Path)
	err = cliCommand.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func itemFunc(repos []Item) func(i int) string {
	return func(i int) string {
		return repos[i].DisplayName()
	}
}

func preview(repos []Item) func(i int, w int, h int) string {
	return func(i, w, h int) string {
		if i == -1 {
			return ""
		}
		return fmt.Sprintf("opening: %s", repos[i].Path)
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
