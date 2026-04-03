package main

import (
	"encoding/json"
	"os"
)

// ExtractionConfig defines the structure of your config file
type ExtractionConfig struct {
	Fields map[string]interface{} `json:"fields"`
}

// LoadConfig reads a JSON file and unmarshals it into ExtractionConfig
func LoadConfig(filename string) (ExtractionConfig, error) {
	var cfg ExtractionConfig

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	// Unmarshal JSON
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
