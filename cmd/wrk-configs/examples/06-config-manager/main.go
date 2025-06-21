package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// ConfigReader универсальный читатель конфигураций
type ConfigReader struct {
	Data map[string]interface{}
}

// NewConfigReader создает новый ConfigReader
func NewConfigReader() *ConfigReader {
	return &ConfigReader{
		Data: make(map[string]interface{}),
	}
}

// ReadJSON читает JSON файл
func (cr *ConfigReader) ReadJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %s: %w", filePath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл %s: %w", filePath, err)
	}

	if err := json.Unmarshal(data, &cr.Data); err != nil {
		return fmt.Errorf("не удалось распарсить JSON из %s: %w", filePath, err)
	}

	return nil
}

// Get получает значение по ключу с поддержкой вложенных ключей
func (cr *ConfigReader) Get(key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	current := cr.Data

	for i, k := range keys {
		if i == len(keys)-1 {
			value, exists := current[k]
			return value, exists
		}

		if next, exists := current[k]; exists {
			if nextMap, ok := next.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	return nil, false
}

// GetAllKeys возвращает все ключи (включая вложенные)
func (cr *ConfigReader) GetAllKeys() []string {
	var keys []string
	cr.collectKeys("", cr.Data, &keys)
	return keys
}

// collectKeys рекурсивно собирает все ключи
func (cr *ConfigReader) collectKeys(prefix string, data map[string]interface{}, keys *[]string) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		*keys = append(*keys, fullKey)

		if nested, ok := value.(map[string]interface{}); ok {
			cr.collectKeys(fullKey, nested, keys)
		}
	}
}

// PrintStructure выводит структуру конфигурации
func (cr *ConfigReader) PrintStructure() {
	fmt.Println("Структура конфигурации:")
	cr.printValue("", cr.Data, 0)
}

// printValue рекурсивно выводит значения
func (cr *ConfigReader) printValue(key string, value interface{}, indent int) {
	indentStr := strings.Repeat("  ", indent)

	switch v := value.(type) {
	case map[string]interface{}:
		if key != "" {
			fmt.Printf("%s%s: {\n", indentStr, key)
		}
		for k, val := range v {
			cr.printValue(k, val, indent+1)
		}
		if key != "" {
			fmt.Printf("%s}\n", indentStr)
		}
	case []interface{}:
		fmt.Printf("%s%s: [\n", indentStr, key)
		for i, item := range v {
			cr.printValue(fmt.Sprintf("[%d]", i), item, indent+1)
		}
		fmt.Printf("%s]\n", indentStr)
	default:
		fmt.Printf("%s%s: %v (%T)\n", indentStr, key, v, v)
	}
}

// ConfigManager управляет различными типами конфигурационных файлов
type ConfigManager struct {
	JSONReader *ConfigReader
	INIData    interface{}
	FilePath   string
	FileType   string
}

// NewConfigManager создает новый менеджер конфигураций
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		JSONReader: NewConfigReader(),
	}
}

// LoadConfig загружает конфигурацию из файла
func (cm *ConfigManager) LoadConfig(filePath string) error {
	cm.FilePath = filePath
	ext := strings.ToLower(filepath.Ext(filePath))
	cm.FileType = ext

	switch ext {
	case ".json":
		return cm.loadJSON()
	case ".ini":
		return cm.loadINI()
	case ".yaml", ".yml":
		return cm.loadYAML()
	default:
		return fmt.Errorf("неподдерживаемый тип файла: %s", ext)
	}
}

// loadJSON загружает JSON конфигурацию
func (cm *ConfigManager) loadJSON() error {
	return cm.JSONReader.ReadJSON(cm.FilePath)
}

// loadINI загружает INI конфигурацию
func (cm *ConfigManager) loadINI() error {
	cfg, err := ini.Load(cm.FilePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения INI файла: %w", err)
	}

	// Преобразуем в map[string]map[string]string
	iniData := make(map[string]map[string]string)

	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		if sectionName == "DEFAULT" {
			sectionName = "default"
		}

		iniData[sectionName] = make(map[string]string)
		for _, key := range section.Keys() {
			iniData[sectionName][key.Name()] = key.String()
		}
	}

	cm.INIData = iniData
	return nil
}

// loadYAML загружает YAML конфигурацию (заглушка)
func (cm *ConfigManager) loadYAML() error {
	return fmt.Errorf("поддержка YAML пока не реализована")
}

// PrintInfo выводит информацию о загруженной конфигурации
func (cm *ConfigManager) PrintInfo() {
	fmt.Printf("Файл: %s\n", cm.FilePath)
	fmt.Printf("Тип: %s\n", cm.FileType)
	fmt.Println(strings.Repeat("-", 40))

	switch cm.FileType {
	case ".json":
		cm.JSONReader.PrintStructure()
	case ".ini":
		fmt.Println("INI конфигурация:")
		if data, ok := cm.INIData.(map[string]map[string]string); ok {
			for section, values := range data {
				fmt.Printf(" [%s]\n", section)
				for key, value := range values {
					fmt.Printf("   %s = %s\n", key, value)
				}
			}
		}
		/*
			case ".ini":
				fmt.Println("INI конфигурация:")
				if data, ok := cm.INIData.(struct{ Section map[string]map[string]string }); ok {
					for section, values := range data.Section {
						fmt.Printf("  [%s]\n", section)
						for key, value := range values {
							fmt.Printf("    %s = %s\n", key, value)
						}
					}
				}
		*/

	}
}

// findConfigFiles ищет конфигурационные файлы в директории
func findConfigFiles(dir string) ([]string, error) {
	var configFiles []string
	supportedExts := []string{".ini", ".json", ".yaml", ".yml"}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать директорию %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))

		for _, supportedExt := range supportedExts {
			if ext == supportedExt {
				configFiles = append(configFiles, filepath.Join(dir, fileName))
				break
			}
		}
	}

	return configFiles, nil
}

const DefaultConfigDir = "cmd/wrk-configs/configs/examples/"

func main() {
	var filesToProcess []string

	if len(os.Args) > 1 {
		// Если указан файл в аргументах, обрабатываем его
		filesToProcess = os.Args[1:]
	} else {
		// Иначе ищем все конфигурационные файлы в текущей директории
		configFiles, err := findConfigFiles(DefaultConfigDir)
		if err != nil {
			fmt.Printf("Ошибка поиска файлов: %v\n", err)
			return
		}
		filesToProcess = configFiles
	}

	if len(filesToProcess) == 0 {
		fmt.Println("Конфигурационные файлы не найдены")
		fmt.Println("Использование: go run config_manager.go [файл1] [файл2] ...")
		return
	}

	fmt.Printf("Найдено %d конфигурационных файлов:\n", len(filesToProcess))
	for i, file := range filesToProcess {
		fmt.Printf("%d. %s\n", i+1, file)
	}

	fmt.Println(strings.Repeat("=", 60))

	// Обрабатываем каждый файл
	manager := NewConfigManager()

	for _, filePath := range filesToProcess {
		fmt.Printf("\nОбработка файла: %s\n", filePath)
		fmt.Println(strings.Repeat("=", 60))

		// Проверяем существование файла
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Файл %s не существует, пропускаем\n", filePath)
			continue
		}

		if err := manager.LoadConfig(filePath); err != nil {
			fmt.Printf("Ошибка загрузки %s: %v\n", filePath, err)
			continue
		}

		manager.PrintInfo()

		// Для JSON файлов показываем дополнительную информацию
		if strings.ToLower(filepath.Ext(filePath)) == ".json" {
			fmt.Println("\nДоступные ключи:")
			keys := manager.JSONReader.GetAllKeys()
			for _, key := range keys {
				if value, exists := manager.JSONReader.Get(key); exists {
					fmt.Printf("  %s: %T\n", key, value)
				}
			}
		}
	}
}
