package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mongo struct {
		URI      string `yaml:"uri"`
		Database string `yaml:"database"`
	} `yaml:"mongo"`
	MySQL struct {
		URI string `yaml:"uri"`
	} `yaml:"mysql"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

var AppConfig Config

func LoadConfig(configPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
	}
}
