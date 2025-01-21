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
	if err != nil {
		return config, err
	}

	for serviceName, service := range config.Services {
		service.Name = serviceName
		config.Services[serviceName] = service
	}

	return config, nil
}
