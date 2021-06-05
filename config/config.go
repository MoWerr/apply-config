package config

import (
	"encoding/json"
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
		return nil, err
	}

	return Read(json, defaultDataSources)
}
