package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

func createLaunchCommands(configs []ConfigItem) []LaunchCommand {
	var launchCommands []LaunchCommand
	resultChannel := make(chan []LaunchCommand, len(configs))
	var wg sync.WaitGroup
	for _, config := range configs {
		wg.Add(1)
		go parseLaunchCommandsFromConfigItem(config, &wg, resultChannel)
	}
	wg.Wait()
	for i := 0; i < len(configs); i++ {
		launchCommands = append(launchCommands, <-resultChannel...)
	}
	return launchCommands
}

func parseLaunchCommandsFromConfigItem(config ConfigItem, wg *sync.WaitGroup, channel chan []LaunchCommand) {
	defer wg.Done()
	listCommand := exec.Command("bash", "-c", config.List)
	output, err := listCommand.Output()
	if err != nil {
		log.Fatal(err)
	}
	var launchCommands []LaunchCommand
	for _, line := range strings.Split(strings.TrimSuffix(string(output), "\n"), "\n") {
		parts := strings.Split(line, ",")
		launchCommand := LaunchCommand{
			Name:    parts[0],
			Command: config.Command,
			Path:    parts[1],
		}
		launchCommands = append(launchCommands, launchCommand)
	}
	channel <- launchCommands
}

type LaunchCommand struct {
	Name    string
	Command string
	Path    string
}

func (launchCommand *LaunchCommand) DisplayName() string {
	return launchCommand.Command + " " + launchCommand.Name
}

func launchCommandsFunc(launchCommands []LaunchCommand) func(index int) string {
	return func(index int) string {
		return launchCommands[index].DisplayName()
	}
}

func previewFunc(launchCommands []LaunchCommand) func(i int, w int, h int) string {
	return func(i, w, h int) string {
		if i == -1 {
			return ""
		}
		return fmt.Sprintf("opening: %s", launchCommands[i].Path)
	}
}
