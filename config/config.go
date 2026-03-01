package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Theme struct {
	AccentColor string `yaml:"accent_color"`
	DimColor    string `yaml:"dim_color"`
}

type Item struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"` // optional subtitle shown below selected item
	Type        string `yaml:"type"`                  // directory | command | url
	Path        string `yaml:"path"`                  // for type=directory
	Command     string `yaml:"command"`               // for type=command
	URL         string `yaml:"url"`                   // for type=url
}

type Category struct {
	Name  string `yaml:"name"`
	Icon  string `yaml:"icon"`
	Color string `yaml:"color,omitempty"` // optional hex accent, e.g. "#ff8800"
	Items []Item `yaml:"items"`
}

type Config struct {
	Theme      Theme      `yaml:"theme"`
	Editor     string     `yaml:"editor"`     // preferred editor binary; falls back to $EDITOR then auto-detect
	Categories []Category `yaml:"categories"`
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pspterm")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func ExampleConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml.example")
}

func Load() (Config, error) {
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func Save(cfg Config) error {
	if err := os.MkdirAll(ConfigDir(), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0644)
}
