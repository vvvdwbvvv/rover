package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

func ParseYAML(filePath string) (RoverCompose, error) {
	var config RoverCompose
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}
