package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func GetConfig() (Config, error) {
	f, err := os.ReadFile("config.yaml")
	c := Config{}
	if err != nil {
		return c, fmt.Errorf("error reading config file: %v", err)

	}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return c, fmt.Errorf("error unmarshalling config file: %v", err)
	}
	return c, nil
}
