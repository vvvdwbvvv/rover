package config

import (
	"errors"
	"fmt"
	"os"
)

func LoadConfig() (RoverCompose, error) {
	var config RoverCompose
	var err error

	if fileExists("rover-compose.yaml") {
		config, err = ParseYAML("rover-compose.yaml")
	} else if fileExists("rover-compose.toml") {
		config, err = ParseTOML("rover-compose.toml")
	} else if fileExists("rover-compose.json") {
		config, err = ParseJSON("rover-compose.json")
	} else {
		return config, errors.New("No rover-compose file found")
	}

	if err != nil {
		return config, fmt.Errorf("Failed to parse config: %v", err)
	}
	return config, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
