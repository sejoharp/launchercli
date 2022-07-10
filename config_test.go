package main

import (
	"reflect"
	"testing"
)

func TestReadAndParseConfigFile(t *testing.T) {
	config := readAndParseConfigFile("test-resources/valid_config.json")
	expected := []ConfigItem{
		{
			DisplayName: "code {0}",
			List:        "for directory in $(ls -d /home/user/repos/*); do echo $(basename $directory),$directory;done",
			Command:     "code {1}",
		},
		{
			DisplayName: "actions {0}",
			List:        "for directory in $(ls -d /home/user/repos/*); do echo $(basename $directory);done",
			Command:     "open https://github.com/company/{0}/actions",
		}}
	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Configs are not equal. Got %v, wanted %v", config, expected)
	}
}
