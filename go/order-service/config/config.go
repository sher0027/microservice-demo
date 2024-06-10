package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MySQL struct {
		URI string `yaml:"uri"`
	} `yaml:"mysql"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
	} `yaml:"kafka"`
	Service struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
	} `yaml:"service"`
	Registry struct {
		URL string `yaml:"url"`
	} `yaml:"registry"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile("config/config.yaml")
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

func RegisterService(cfg *Config) error {
	service := map[string]string{
		"name": cfg.Service.Name,
		"addr": fmt.Sprintf("localhost:%d", cfg.Service.Port),
	}
	data, err := json.Marshal(service)
	if err != nil {
		return err
	}

	resp, err := http.Post(cfg.Registry.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service: %v", resp.Status)
	}

	return nil
}
