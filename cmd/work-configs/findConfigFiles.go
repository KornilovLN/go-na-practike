package main

import (
	"encoding/json"
	"fmt"

	//"github.com/KornilovLN/go-na-praktike/author"

	"os"
	"path/filepath"
	"strings"

	"github.com/KornilovLN/go-na-practike/cmd/work-configs/author"
	"gopkg.in/gcfg.v1" // Импорт пакета для работы с INI-файлами
	"gopkg.in/yaml.v3" // Импорт пакета для работы с YAML-файлами
)

// findConfigFiles ищет конфигурационные файлы в указанной директории и ее поддиректориях
func findConfigFiles(dir string) ([]string, error) {
	var configFiles []string

	// Поддерживаемые расширения конфигурационных файлов
	supportedExts := []string{".ini", ".json", ".yaml", ".yml"}

	// Читаем содержимое директории
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("err: findConfigFiles: не удалось прочитать директорию %s: %w", dir, err)
	}

	// каждый файл - если это файл, то проверяем его расширение
	for _, entry := range entries {
		if entry.IsDir() { // Пропускаем директории
			continue
		}

		// Получаем имя файла
		fileName := entry.Name()

		// Получаем расширение в нижнем регистре
		ext := strings.ToLower(filepath.Ext(fileName))

		// является ли расширение поддерживаемым согласно заданию?
		for _, supportedExt := range supportedExts {
			if ext == supportedExt {
				configFiles = append(configFiles, filepath.Join(dir, fileName))
				break
			}
		}
	}

	return configFiles, nil
}

// readINIConfig читает INI файл
func readINIConfig(filePath string) error {
	config := struct {
		Section struct {
			Enabled bool
			Path    string
		}
	}{}

	err := gcfg.ReadFileInto(&config, filePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения INI файла %s: %w", filePath, err)
	}

	fmt.Printf("INI файл: %s\n", filePath)
	fmt.Printf("\tEnabled: %t\n", config.Section.Enabled)
	fmt.Printf("\tPath: %s\n", config.Section.Path)

	return nil
}

// readJSONConfig читает JSON файл
func readJSONConfig(filePath string) error {
	// Определение структуры под JSON-файл
	type configuration struct {
		Enabled bool
		Path    string
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("err: readJSONConfig: ошибка открытия JSON файла %s: %w", filePath, err)
	}
	// defer для закрытия файла после завершения работы с ним
	defer func() {
		if err3 := file.Close(); err3 != nil {
			fmt.Printf("err: readJSONConfig: Ошибка при закрытии файла: %v", err3)
		}
	}()

	// Создание нового декодера JSON
	decoder := json.NewDecoder(file)

	// Извлечение JSON-значений в переменные
	conf := configuration{}
	err2 := decoder.Decode(&conf)
	if err2 != nil {
		fmt.Println("--- Error:", err2)
	}

	// Вывод значений полей json на экран
	fmt.Println("JSON файл: conf.json:")
	fmt.Println("\tEnabled: ", conf.Enabled)
	fmt.Println("\tPath: ", conf.Path)

	// или так:
	//fmt.Println("\nФорматированный вывод conf.json:")
	//fmt.Printf("{\n\tEnabled: %t,\n\tPath: %s\n}", conf.Enabled, conf.Path)

	return nil
}

// readYAMLConfig читает YAML файл
func readYAMLConfig(filePath string) error {
	type Config struct {
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	}

	// Чтение YAML файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("err: readYAMLConfig: Ошибка чтения файла: ", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Println("err: readYAMLConfig: Ошибка парсинга YAML: ", err)
	}

	fmt.Println("YAML файл: conf.yaml:")
	fmt.Println("\tEnabled: ", config.Enabled)
	fmt.Println("\tPath: ", config.Path)

	return nil
}

func main() {
	authorInfo := author.NewAuthorInfo()
	authorInfo.Print()

	// Указываем директорию для поиска (можно изменить на нужную)
	searchDir := "."

	// Если передан аргумент командной строки, используем его как директорию
	if len(os.Args) > 1 {
		searchDir = os.Args[1]
	}

	// Ищем конфигурационные файлы
	configFiles, err := findConfigFiles(searchDir)
	if err != nil {
		fmt.Printf("err: main: Ошибка поиска файлов: %v\n", err)
		return
	}

	if len(configFiles) == 0 {
		fmt.Println("main: Конфигурационные файлы не найдены")
		return
	} else {
		fmt.Printf("main: Найдено %d конфигурационных файлов:\n\n", len(configFiles))
		for i, file := range configFiles {
			fmt.Printf("%d. %s\n", i+1, file)
		}
	}

	fmt.Println("\nmain: Обработка найденных файлов:\n")

	// Обрабатываем каждый найденный файл
	for i, filePath := range configFiles {
		//fmt.Printf("%d. --- %s ---\n", i+1, filePath)
		fmt.Printf("%d. ", i+1)

		// Получаем расширение файла и выполняем соответствующую обработку
		ext := strings.ToLower(filepath.Ext(filePath))
		switch ext {
		case ".ini":
			if err := readINIConfig(filePath); err != nil {
				fmt.Printf("err: main: Ошибка обработки %s: %v\n", filePath, err)
			}
		case ".json":
			if err := readJSONConfig(filePath); err != nil {
				readJSONConfig(filePath)
			}
		case ".yaml", ".yml":
			if err := readYAMLConfig(filePath); err != nil {
				readYAMLConfig(filePath)
			}
		default:
			fmt.Printf("err: main: Неверное расширение %s: %v\n", filePath, err)
		}
		fmt.Println()
	}
}
