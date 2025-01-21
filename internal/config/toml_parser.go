package config

import (
	"github.com/BurntSushi/toml"
)

func ParseTOML(filePath string) (RoverCompose, error) {
	var config RoverCompose
	_, err := toml.DecodeFile(filePath, &config)
	return config, err
}
