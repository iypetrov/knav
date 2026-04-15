package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Current string   `yaml:"current"`
	Targets []Target `yaml:"targets"`
}

func (cfg *Config) Save() error {
	var configPath string

	if os.Getenv("KNAVCONFIG") == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath = filepath.Join(homeDir, ".config", "knav", "config.yaml")
	} else {
		configPath = filepath.Join(os.Getenv("KNAVCONFIG"), "config.yaml")
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (cfg *Config) PickTarget() (Target, error) {
	cmd := exec.Command("fzf", "--ansi", "--no-preview")

	var stdin bytes.Buffer
	for _, t := range cfg.Targets {
		stdin.WriteString(t.Name + "\n")
	}
	cmd.Stdin = &stdin

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return Target{}, err
	}

	selected := strings.TrimSpace(stdout.String())

	for _, t := range cfg.Targets {
		if t.Name == selected {
			cfg.Current = t.Name
			if err := cfg.Save(); err != nil {
				return Target{}, err
			}

			return t, nil
		}
	}

	return Target{}, errors.New("no target selected")
}

func (cfg Config) CurrentTarget() (Target, error) {
	for _, t := range cfg.Targets {
		if t.Name == cfg.Current {
			return t, nil
		}
	}
	return Target{}, errors.New("current target does not exist in the target list")
}

type Target struct {
	Name           string   `yaml:"name"`
	KubeconfigPath string   `yaml:"kubeconfigPath"`
	Envs           []Env    `yaml:"envs,omitempty"`
	Restricted     bool     `yaml:"restricted"`
	AllowedActions []string `yaml:"allowedActions,omitempty"`
}

type Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func LoadConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	if os.Getenv("KNAVCONFIG") == "" {
		viper.AddConfigPath(filepath.Join(home, ".config", "knav"))
	} else {
		viper.AddConfigPath(os.Getenv("KNAVCONFIG"))
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	return cfg, err
}

func DefaultConfig() Config {
	return Config{
		Current: "local",
		Targets: []Target{
			{
				Name:           "local",
				KubeconfigPath: "~/.kube/config",
				Restricted:     false,
			},
		},
	}
}

func CreateDefaultConfigIfEmpty() error {
	var configPath string

	if os.Getenv("KNAVCONFIG") == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath = filepath.Join(homeDir, ".config", "knav", "config.yaml")
	} else {
		configPath = filepath.Join(os.Getenv("KNAVCONFIG"), "config.yaml")
	}

	if _, err := os.Stat(configPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	cfg := DefaultConfig()

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
