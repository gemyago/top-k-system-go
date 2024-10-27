package config

import (
	"embed"
	"fmt"

	"github.com/spf13/viper"
)

//go:embed *.json
var resources embed.FS

func mergeResourceCfg(cfg *viper.Viper, resourceName string) error {
	resourceStream, err := resources.Open(resourceName)
	if err != nil {
		return fmt.Errorf("failed to read config %v: %w", resourceName, err)
	}
	defer resourceStream.Close()

	if err = cfg.MergeConfig(resourceStream); err != nil {
		return fmt.Errorf("failed to load config %v: %w", resourceName, err)
	}
	return nil
}

type LoadOpts struct {
	env                   string
	defaultConfigFileName string
}

func (opts *LoadOpts) WithEnv(val string) *LoadOpts {
	if val != "" {
		opts.env = val
	}
	return opts
}

func NewLoadOpts() *LoadOpts {
	return &LoadOpts{
		env:                   "local",
		defaultConfigFileName: "default.json",
	}
}

func New() *viper.Viper {
	v := viper.New()
	v.SetConfigType("json")
	return v
}

func Load(cfg *viper.Viper, opts *LoadOpts) error {
	if err := mergeResourceCfg(cfg, opts.defaultConfigFileName); err != nil {
		return err
	}

	if err := mergeResourceCfg(cfg, opts.env+".json"); err != nil {
		return err
	}

	return nil
}
