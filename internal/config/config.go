package config

import (
	"fmt"
	"go-heal/internal/types"
	"os"

	"gopkg.in/yaml.v2"
)

func Load(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Config file could not be read: %w", err)
	}

	var cfg types.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("yaml parse error: %w", err)
	}

	return &cfg, nil
}