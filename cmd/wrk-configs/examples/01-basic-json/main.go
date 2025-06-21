// main.go
package main

import (
	"fmt"
	"log"

	"github.com/KornilovLN/go-na-practike/cmd/wrk-configs/pkg/parsers"
	"github.com/KornilovLN/go-na-practike/cmd/wrk-configs/pkg/types"
)

func main() {
	fmt.Println("=== Пример 1: Базовая работа с JSON ===")

	parser := parsers.NewJSONParser()

	var config types.CommonConfig
	err := parser.ParseFile("cmd/wrk-configs/configs/examples/app.json", &config)
	if err != nil {
		log.Fatal("Ошибка парсинга:", err)
	}

	fmt.Printf("Конфигурация загружена:\n")
	fmt.Printf("  Database: %s:%d (user: %s)\n",
		config.Database.Host, config.Database.Port, config.Database.Username)
	fmt.Printf("  Server: %s:%d\n",
		config.Server.Host, config.Server.Port)
	fmt.Printf("  Debug: %v\n", config.Debug)
	fmt.Printf("  Logging: %s -> %s\n",
		config.Logging.Level, config.Logging.File)
}
