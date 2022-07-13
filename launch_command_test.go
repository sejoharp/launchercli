package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPartsToLaunchCommand(t *testing.T) {
	got := toLaunchCommand(ConfigItem{
		DisplayName: "code {0}",
		List:        "",
		Command:     "code {1}",
	}, "repo,/path/to/repo")
	expected := LaunchCommand{
		DisplayName: "code repo",
		Command:     "code /path/to/repo",
	}
	if got != expected {
		t.Errorf("Header is not equal. Got %v, wanted %v", got, expected)
	}
}

func TestReplacePlaceholder(t *testing.T) {
	parts := []string{"part1", "part2"}
	stringWithPlaceholder := "code {1}"

	result := replacePlaceholder(stringWithPlaceholder, parts)

	assert.Equal(t, "code part2", result)
}

func TestReturnOriginalCommandIfNotPlaceholderIsPresent(t *testing.T) {
	parts := []string{"part1", "part2"}
	stringWithPlaceholder := "code"

	result := replacePlaceholder(stringWithPlaceholder, parts)

	assert.Equal(t, "code", result)
}

func TestReturnOriginalCommandIfPartIdIsNotPresent(t *testing.T) {
	parts := []string{"part1", "part2"}
	stringWithPlaceholder := "code {2}"

	result := replacePlaceholder(stringWithPlaceholder, parts)

	assert.Equal(t, "code {2}", result)
}

func TestReplaceMultiplePlaceholder(t *testing.T) {
	parts := []string{"-m", "/tmp/myproject"}
	stringWithPlaceholder := "code {0} {1}"

	result := replacePlaceholder(stringWithPlaceholder, parts)

	assert.Equal(t, "code -m /tmp/myproject", result)
}
