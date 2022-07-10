package main

import "testing"

func TestPartsToLaunchCommand(t *testing.T) {
	got := partsToLaunchCommand(ConfigItem{
		DisplayName: "code {0}",
		List:        "",
		Command:     "code {1}",
	}, []string{"repo", "/path/to/repo"})
	expected := LaunchCommand{
		DisplayName: "code repo",
		Command:     "code /path/to/repo",
	}
	if got != expected {
		t.Errorf("Header is not equal. Got %v, wanted %v", got, expected)
	}
}
