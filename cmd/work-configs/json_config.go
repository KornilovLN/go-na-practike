package main

// Анализ конфигурационного JSON-файла: json_config.go

import (
	"encoding/json"
	"fmt"
	"os"
)

// Определение структуры под JSON-файл
type configuration struct {
	Enabled bool
	Path    string
}

func main() {
	// Открытие конфигурационного файла
	//file, err := os.Open("cmd/work-configs/conf.json")
	file, err := os.Open("conf.json")
	if err != nil {
		fmt.Println("--- Error opening file:", err)
		return
	}
	// Использование defer для закрытия файла после завершения работы с ним
	defer func() {
		if err3 := file.Close(); err3 != nil {
			fmt.Printf("--- Ошибка при закрытии файла: %v", err3)
		}
	}()
	// Альтернативный способ закрытия файла без проверки ошибки:
	// defer file.Close()

	// Создание нового декодера JSON
	decoder := json.NewDecoder(file)

	// Извлечение JSON-значений в переменные
	conf := configuration{}
	err2 := decoder.Decode(&conf)
	if err2 != nil {
		fmt.Println("--- Error:", err2)
	}

	// Вывод значений полей json на экран
	fmt.Println("Конфигурация conf.json:")
	fmt.Println(conf.Enabled)
	fmt.Println(conf.Path)

	// или так:
	fmt.Println("\nФорматированный вывод conf.json:")
	fmt.Printf("{\n\tEnabled: %t,\n\tPath: %s\n}", conf.Enabled, conf.Path)
}
