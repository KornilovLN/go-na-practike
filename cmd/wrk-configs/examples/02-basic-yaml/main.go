package main

// Анализ конфигурационного YAML-файла: yaml_config2.go
// доустановить лучший пакет gopkg.in/yaml.v3
// go get gopkg.in/yaml.v3

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

func main() {
	//data, err := os.ReadFile("cmd/work-configs/conf.yaml")
	data, err := os.ReadFile("../../configs/examples/conf.yaml")
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Ошибка парсинга YAML: %v", err)
	}

	fmt.Printf("Enabled: %v\nPath: %v\n", config.Enabled, config.Path)
}
