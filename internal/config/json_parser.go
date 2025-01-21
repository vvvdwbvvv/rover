package config

import (
	"encoding/json"
	"os"
)

func ParseJSON(filePath string) (RoverCompose, error) {
	var config RoverCompose
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}
