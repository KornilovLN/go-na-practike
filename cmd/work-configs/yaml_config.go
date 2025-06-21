package main

// Анализ конфигурационного YAML-файла: yaml_config.go
// Надо добавить в go.mod:
// go get github.com/kylelemons/go-gypsy/yaml
// Импортирование необходимых пакетов

import (
	"fmt"
	"log"

	"github.com/kylelemons/go-gypsy/yaml" // Импорт пакета YAML сторонних
)

func main() {
	// Чтение и анализ YAML-файла
	//config, err := yaml.ReadFile("cmd/work-configs/conf.yaml")
	config, err := yaml.ReadFile("conf.yaml")
	if err != nil {
		log.Fatalf("Ошибка чтения файла конфигурации: %v", err)
	}

	// Получаем значения с проверкой ошибок
	path, err := config.Get("path")
	if err != nil {
		log.Printf("Ошибка получения path: %v", err)
	} else {
		fmt.Println("Path:", path)
	}

	enabled, err := config.GetBool("enabled")
	if err != nil {
		log.Printf("Ошибка получения enabled: %v", err)
	} else {
		fmt.Println("Enabled:", enabled)
	}
}
