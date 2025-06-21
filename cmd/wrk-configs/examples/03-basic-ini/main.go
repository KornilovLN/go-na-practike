package main

// Извлечение данных из конфигурационного INI-файла: ini_config.go
// Для работы с INI-файлами используется пакет gcfg
// go get gopkg.in/gcfg.v1

import (
	"fmt"

	"gopkg.in/gcfg.v1" // Импорт пакета для работы с INI-файлами
)

func main() {
	config := struct { // Создание структуры для конфигурационных значений
		Section struct {
			Enabled bool
			Path    string
		}
	}{}

	// Извлечение данных из INI-файла в структуру с обработкой ошибок
	//err := gcfg.ReadFileInto(&config, "cmd/work-configs/conf.ini")
	err := gcfg.ReadFileInto(&config, "../../configs/examples/conf.ini")
	if err != nil {
		fmt.Println("Failed to parse config file: %s", err)
	}

	fmt.Println("Использование значений из INI-файла: ")
	fmt.Println(config.Section.Enabled)
	fmt.Println(config.Section.Path)
}
