package config

import (
    "gopkg.in/yaml.v2"
    "os"
)

type Config struct {
    MySQL struct {
        URI string `yaml:"uri"`
    } `yaml:"mysql"`
    Server struct {
        Port int `yaml:"port"`
    } `yaml:"server"`
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