package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Kafka struct {
		Brokers []string `yaml:"brokers"`
	} `yaml:"kafka"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
