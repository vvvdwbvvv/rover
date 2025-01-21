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
	if err != nil {
		return config, err
	}

	for serviceName, service := range config.Services {
		service.Name = serviceName
		config.Services[serviceName] = service
	}

	return config, nil
}
