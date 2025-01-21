package config

type Service struct {
	Name        string            `yaml:"name" toml:"name" json:"name"`
	Image       string            `yaml:"image" toml:"image" json:"image"`
	Command     []string          `yaml:"command,omitempty" toml:"command,omitempty" json:"command,omitempty"`
	Ports       []string          `yaml:"ports,omitempty" toml:"ports,omitempty" json:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty" toml:"volumes,omitempty" json:"volumes,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty" toml:"environment,omitempty" json:"environment,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty" toml:"depends_on,omitempty" json:"depends_on,omitempty"`
}

type RoverCompose struct {
	Version  string             `yaml:"version" toml:"version" json:"version"`
	Services map[string]Service `yaml:"services" toml:"services" json:"services"`
	Volumes  map[string]struct {
		Driver string `yaml:"driver,omitempty" toml:"driver,omitempty" json:"driver,omitempty"`
	} `yaml:"volumes,omitempty" toml:"volumes,omitempty" json:"volumes,omitempty"`
}
