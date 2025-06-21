package main

// wrk-config-creator.go
// Package main creates a directory structure based on a JSON file.
// usage: go run create_dirs.go
//		  go run create_dirs.go my-structure.json
//		  go run create_dirs.go /path/to/config.json

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Structure представляет структуру JSON файла с описанием директорий для создания
// Base - базовый путь, относительно которого создаются все директории
// Structure - вложенная структура директорий в виде map[string]interface{}
type Structure struct {
	Base      string                 `json:"base"`      // Корн. дир. для создания структуры
	Structure map[string]interface{} `json:"structure"` // Иерархическое описание поддиректорий
}

// createDirs рекурсивно создает директории на основе переданной структуры
// basePath - текущий базовый путь для создания директорий
// structure - map с описанием директорий,
//
//	где ключ - имя директории,
//	значение - поддиректории
//
// Возвращает ошибку если не удалось создать какую-либо директорию
func createDirs(basePath string, structure map[string]interface{}) error {
	// Итерируемся по всем элементам структуры
	for name, children := range structure {
		// currentPath - полный путь к создаваемой директории
		currentPath := filepath.Join(basePath, name)
		fmt.Printf("Creating: %s\n", currentPath)

		// Создаем директорию с правами 0755
		if err := os.MkdirAll(currentPath, 0755); err != nil {
			return err
		}

		// Проверяем есть ли поддиректории и рекурсивно их создаем
		// childMap - приведение interface{} к map[string]interface{} для поддиректорий
		if childMap, ok := children.(map[string]interface{}); ok && len(childMap) > 0 {
			if err := createDirs(currentPath, childMap); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	// jsonFile - путь к JSON файлу с описанием структуры директорий
	// По умолчанию "structure.json",
	// но может быть переопределен через аргумент командной строки
	jsonFile := "wrk-config-creator.json"
	if len(os.Args) > 1 {
		jsonFile = os.Args[1]
	}

	// s - структура для хранения распарсенных данных из JSON
	var s Structure

	// data - содержимое JSON файла в виде байтов
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// Парсим JSON в структуру s
	if err := json.Unmarshal(data, &s); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Создаем директории начиная с базового пути
	if err := createDirs(s.Base, s.Structure); err != nil {
		log.Fatal("Error creating dirs:", err)
	}

	fmt.Println("Done!")
}
