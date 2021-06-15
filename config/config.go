package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type File struct {
	Source      string
	Destination string
}

type Instance struct {
	Name        string
	Destination string
}

type Config struct {
	Files       []File
	Instances   []Instance
	DataSources []string
}

func Read(configJSON []byte, defaultDataSources []string) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(configJSON, config)
	if err != nil {
		return nil, err
	}

	if len(config.DataSources) == 0 {
		config.DataSources = defaultDataSources
	}

	return config, err
}

func ReadFile(filename string, defaultDataSources []string) (*Config, error) {
	json, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the config file: %q; %w", filename, err)
	}

	config, err := Read(json, defaultDataSources)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the config json: %q; %w", filename, err)
	}

	return config, nil
}

func (r Config) GetTemplatePaths() []string {
	t := make([]string, 0, len(r.Files))
	for _, v := range r.Files {
		t = append(t, v.Source)
	}
	return t
}
