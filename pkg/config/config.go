package config

import (
	"fmt"
	"os"

	"netradio/pkg/cloud"
	"netradio/pkg/database"
	"netradio/pkg/email"
	"netradio/pkg/jwt"

	"gopkg.in/yaml.v2"
)

const (
	DefaultYAMLPath = "config.yaml"
)

type Config struct {
	Port      int             `yaml:"port"`
	Jwt       jwt.Config      `yaml:"jwt"`
	Database  database.Config `yaml:"database"`
	Email     email.Config    `yaml:"email"`
	Cloud     cloud.Config    `yaml:"cloud"`
	Streaming bool            `yaml:"streaming"`
}

func NewConfigFromYAML(path string) (Config, error) {
	config := &Config{}
	file, err := os.Open(path)
	fmt.Println()
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return Config{}, err
	}

	return *config, nil
}
