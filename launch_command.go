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

func toLaunchCommand(config ConfigItem, line string) LaunchCommand {
	parts := strings.Split(line, ",")
	displayName := config.DisplayName
	command := config.Command
	return LaunchCommand{
		DisplayName: replacePlaceholder(displayName, parts),
		Command:     replacePlaceholder(command, parts),
	}
}

func replacePlaceholder(placeholderString string, parts []string) string {
	for i := 0; i < len(parts); i++ {
		placeholderString = strings.Replace(placeholderString, fmt.Sprintf("{%d}", i), parts[i], -1)
	}
	return placeholderString
}

func execShellCommand(commandString string) string {
	listCommand := exec.Command("bash", "-c", commandString)
	output, err := listCommand.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(output)
}

// TODO test me
func parseLaunchCommandsFromConfigItem(config ConfigItem, wg *sync.WaitGroup, channel chan []LaunchCommand) {
	defer wg.Done()
	output := execShellCommand(config.List)
	var launchCommands []LaunchCommand
	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		launchCommand := toLaunchCommand(config, line)
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
