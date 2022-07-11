package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

// TODO maybe move bash-command even out of here to test parallelism?
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

// TODO find a better function name
func partsToLaunchCommand(config ConfigItem, parts []string) LaunchCommand {
	displayName := config.DisplayName
	command := config.Command
	// TODO move to separate function and call twice (displayname and command)
	for i := 0; i < len(parts); i++ {
		displayName = strings.Replace(displayName, fmt.Sprintf("{%d}", i), parts[i], -1)
		command = strings.Replace(command, fmt.Sprintf("{%d}", i), parts[i], -1)
	}
	return LaunchCommand{
		DisplayName: displayName,
		Command:     command,
	}
}

func execListCommand(listCommandString string) string {
	listCommand := exec.Command("bash", "-c", listCommandString)
	output, err := listCommand.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(output)
}

// TODO test me
func parseLaunchCommandsFromConfigItem(config ConfigItem, wg *sync.WaitGroup, channel chan []LaunchCommand) {
	defer wg.Done()
	output := execListCommand(config.List)
	var launchCommands []LaunchCommand
	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		parts := strings.Split(line, ",")
		launchCommand := partsToLaunchCommand(config, parts)
		launchCommands = append(launchCommands, launchCommand)
	}
	channel <- launchCommands
}

type LaunchCommand struct {
	DisplayName string
	Command     string
}

// TODO test me
func launchCommandsFunc(launchCommands []LaunchCommand) func(index int) string {
	return func(index int) string {
		return launchCommands[index].DisplayName
	}
}

// TODO test me
func previewFunc(launchCommands []LaunchCommand) func(i int, w int, h int) string {
	return func(i, w, h int) string {
		if i == -1 {
			return ""
		}
		return fmt.Sprintf("opening: %s", launchCommands[i].Command)
	}
}
