package main

import (
	"reflect"
	"testing"
)

func TestReadAndParseConfigFile(t *testing.T) {
	config := readAndParseConfigFile("test-resources/valid_config.json")
	expected := []ConfigItem{{
		List:    "for directory in $(ls -d /home/user/repos/*); do echo $(basename $directory),$directory;done",
		Command: "",
	}}
	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Configs are not equal. Got %v, wanted %v", config, expected)
	}
}
