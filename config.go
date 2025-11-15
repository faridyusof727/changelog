package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GitPath      string       `yaml:"git_path"`
	Ignore       string       `yaml:"ignore"`
	CommitGroups CommitGroups `yaml:"commit_groups"`
}

type CommitGroups struct {
	TitleMaps map[string]string `yaml:"title_maps"`
}

// NewConfig loads and parses the configuration from the given YAML file path.
func NewConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
