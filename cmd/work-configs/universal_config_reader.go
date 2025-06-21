// universal_config_reader.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// Get получает значение по ключу с поддержкой вложенных ключей (например, "database.host")
func (cr *ConfigReader) Get(key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	current := cr.Data

	for i, k := range keys {
		if i == len(keys)-1 {
			// Последний ключ
			value, exists := current[k]
			return value, exists
		}

		// Промежуточный ключ
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

// GetString получает строковое значение
func (cr *ConfigReader) GetString(key string) (string, error) {
	value, exists := cr.Get(key)
	if !exists {
		return "", fmt.Errorf("ключ '%s' не найден", key)
	}

	if str, ok := value.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("значение по ключу '%s' не является строкой", key)
}

// GetInt получает целочисленное значение
func (cr *ConfigReader) GetInt(key string) (int64, error) {
	value, exists := cr.Get(key)
	if !exists {
		return 0, fmt.Errorf("ключ '%s' не найден", key)
	}

	switch v := value.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("значение по ключу '%s' не является числом", key)
	}
}

// GetBool получает булево значение
func (cr *ConfigReader) GetBool(key string) (bool, error) {
	value, exists := cr.Get(key)
	if !exists {
		return false, fmt.Errorf("ключ '%s' не найден", key)
	}

	if b, ok := value.(bool); ok {
		return b, nil
	}

	return false, fmt.Errorf("значение по ключу '%s' не является булевым", key)
}

// GetFloat получает значение с плавающей точкой
func (cr *ConfigReader) GetFloat(key string) (float64, error) {
	value, exists := cr.Get(key)
	if !exists {
		return 0, fmt.Errorf("ключ '%s' не найден", key)
	}

	if f, ok := value.(float64); ok {
		return f, nil
	}

	return 0, fmt.Errorf("значение по ключу '%s' не является числом с плавающей точкой", key)
}

// GetArray получает массив значений
func (cr *ConfigReader) GetArray(key string) ([]interface{}, error) {
	value, exists := cr.Get(key)
	if !exists {
		return nil, fmt.Errorf("ключ '%s' не найден", key)
	}

	if arr, ok := value.([]interface{}); ok {
		return arr, nil
	}

	return nil, fmt.Errorf("значение по ключу '%s' не является массивом", key)
}

// GetStringArray получает массив строк
func (cr *ConfigReader) GetStringArray(key string) ([]string, error) {
	arr, err := cr.GetArray(key)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(arr))
	for i, item := range arr {
		if str, ok := item.(string); ok {
			result[i] = str
		} else {
			return nil, fmt.Errorf("элемент %d массива '%s' не является строкой", i, key)
		}
	}

	return result, nil
}

// GetObject получает вложенный объект
func (cr *ConfigReader) GetObject(key string) (map[string]interface{}, error) {
	value, exists := cr.Get(key)
	if !exists {
		return nil, fmt.Errorf("ключ '%s' не найден", key)
	}

	if obj, ok := value.(map[string]interface{}); ok {
		return obj, nil
	}

	return nil, fmt.Errorf("значение по ключу '%s' не является объектом", key)
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

// ToJSON конвертирует данные обратно в JSON
func (cr *ConfigReader) ToJSON(pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(cr.Data, "", "  ")
	}
	return json.Marshal(cr.Data)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run universal_config_reader.go <путь_к_json_файлу>")
		return
	}

	filePath := os.Args[1]

	// Проверяем существование файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Файл %s не существует\n", filePath)
		return
	}

	// Проверяем расширение файла
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".json" {
		fmt.Printf("Поддерживаются только JSON файлы, получен: %s\n", ext)
		return
	}

	fmt.Printf("Чтение конфигурационного файла: %s\n", filePath)
	fmt.Println(strings.Repeat("=", 50))

	// Создаем ридер и читаем файл
	reader := NewConfigReader()
	if err := reader.ReadJSON(filePath); err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		return
	}

	// Выводим структуру
	reader.PrintStructure()

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Все доступные ключи:")
	keys := reader.GetAllKeys()
	for i, key := range keys {
		fmt.Printf("%d. %s\n", i+1, key)
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Демонстрация получения значений:")

	// Демонстрируем получение различных типов значений
	for _, key := range keys {
		value, exists := reader.Get(key)
		if !exists {
			continue
		}

		switch value.(type) {
		case string:
			if str, err := reader.GetString(key); err == nil {
				fmt.Printf("String - %s: \"%s\"\n", key, str)
			}
		case float64:
			if f, err := reader.GetFloat(key); err == nil {
				fmt.Printf("Float - %s: %g\n", key, f)
			}
			if i, err := reader.GetInt(key); err == nil {
				fmt.Printf("Int - %s: %d\n", key, i)
			}
		case bool:
			if b, err := reader.GetBool(key); err == nil {
				fmt.Printf("Bool - %s: %t\n", key, b)
			}
		case []interface{}:
			if arr, err := reader.GetArray(key); err == nil {
				fmt.Printf("Array - %s: %d элементов\n", key, len(arr))
			}
		case map[string]interface{}:
			if obj, err := reader.GetObject(key); err == nil {
				fmt.Printf("Object - %s: %d полей\n", key, len(obj))
			}
		}
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("JSON представление (отформатированное):")
	if jsonData, err := reader.ToJSON(true); err == nil {
		fmt.Println(string(jsonData))
	}
}
