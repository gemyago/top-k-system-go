package config

import (
	"embed"
	"fmt"

	"github.com/spf13/viper"
)

//go:embed *.json
var resources embed.FS

func Load() (*viper.Viper, error) {
	defaultCfg, err := resources.Open("default.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read default config: %w", err)
	}

	cfg := viper.New()
	cfg.SetConfigType("json")

	if err = cfg.MergeConfig(defaultCfg); err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}
	return cfg, nil
}
