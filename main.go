package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

func homeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return dirname
}

func main() {
	path := homeDir() + "/.launcher-config.json"
	configs := readAndParseConfigFile(path)
	launchCommands := createLaunchCommands(configs)

	index, err := fuzzyfinder.Find(
		launchCommands,
		launchCommandsFunc(launchCommands),
		fuzzyfinder.WithPreviewWindow(previewFunc(launchCommands)))
	if err != nil {
		log.Fatal(err)
	}
	before, after, _ := strings.Cut(launchCommands[index].Command, " ")
	cliCommand := exec.Command(before, after)
	err = cliCommand.Run()
	if err != nil {
		log.Fatal(err)
	}
}
