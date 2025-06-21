# Чтение JSON файлов с неизвестной структурой
    Решения от простого к сложному:

## 1. Использование map[string]interface{} для простых случаев
```GO
//dynamic_json_reader.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// readJSONToMap читает JSON в map[string]interface{}
func readJSONToMap(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s: %w", filePath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл %s: %w", filePath, err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить JSON из %s: %w", filePath, err)
	}

	return result, nil
}

// printJSONStructure выводит структуру JSON с отступами
func printJSONStructure(data interface{}, indent int) {
	indentStr := strings.Repeat("  ", indent)
	
	switch v := data.(type) {
	case map[string]interface{}:
		fmt.Printf("%s{\n", indentStr)
		for key, value := range v {
			fmt.Printf("%s  %s: ", indentStr, key)
			switch value.(type) {
			case map[string]interface{}:
				fmt.Println()
				printJSONStructure(value, indent+2)
			case []interface{}:
				fmt.Printf("[\n")
				if len(value.([]interface{})) > 0 {
					printJSONStructure(value.([]interface{})[0], indent+2)
					if len(value.([]interface{})) > 1 {
						fmt.Printf("%s    ... (%d элементов)\n", indentStr, len(value.([]interface{})))
					}
				}
				fmt.Printf("%s  ]\n", indentStr)
			default:
				fmt.Printf("%v (%s)\n", value, reflect.TypeOf(value))
			}
		}
		fmt.Printf("%s}\n", indentStr)
	case []interface{}:
		fmt.Printf("%s[\n", indentStr)
		for i, item := range v {
			fmt.Printf("%s  [%d]: ", indentStr, i)
			printJSONStructure(item, indent+1)
		}
		fmt.Printf("%s]\n", indentStr)
	default:
		fmt.Printf("%s%v (%s)\n", indentStr, v, reflect.TypeOf(v))
	}
}

func main() {
	filePath := "config.json"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	fmt.Printf("Чтение JSON файла: %s\n", filePath)
	
	data, err := readJSONToMap(filePath)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Println("\nСтруктура JSON:")
	printJSONStructure(data, 0)
	
	fmt.Println("\nДанные:")
	prettyJSON, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(prettyJSON))
}
```

## 2. Генерация Go структур из JSON
```GO
// json_to_struct.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

// StructField представляет поле структуры
type StructField struct {
	Name     string
	Type     string
	JSONTag  string
	Optional bool
}

// StructInfo содержит информацию о структуре
type StructInfo struct {
	Name   string
	Fields []StructField
}

// JSONToStructGenerator генерирует Go структуры из JSON
type JSONToStructGenerator struct {
	structs map[string]*StructInfo
}

// NewJSONToStructGenerator создает новый генератор
func NewJSONToStructGenerator() *JSONToStructGenerator {
	return &JSONToStructGenerator{
		structs: make(map[string]*StructInfo),
	}
}

// toPascalCase преобразует строку в PascalCase
func toPascalCase(s string) string {
	if s == "" {
		return s
	}
	
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	
	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])))
			if len(word) > 1 {
				result.WriteString(strings.ToLower(word[1:]))
			}
		}
	}
	
	if result.Len() == 0 {
		return "Field"
	}
	
	return result.String()
}

// analyzeValue анализирует значение и определяет его тип
func (g *JSONToStructGenerator) analyzeValue(value interface{}, fieldName string) string {
	switch v := value.(type) {
	case nil:
		return "*interface{}"
	case bool:
		return "bool"
	case float64:
		// JSON числа всегда float64, но проверим, целое ли это
		if v == float64(int64(v)) {
			return "int64"
		}
		return "float64"
	case string:
		return "string"
	case []interface{}:
		if len(v) == 0 {
			return "[]interface{}"
		}
		// Анализируем первый элемент массива
		elementType := g.analyzeValue(v[0], fieldName+"Item")
		return "[]" + elementType
	case map[string]interface{}:
		structName := toPascalCase(fieldName)
		if structName == "" {
			structName = "NestedStruct"
		}
		g.analyzeObject(v, structName)
		return structName
	default:
		return "interface{}"
	}
}

// analyzeObject анализирует объект и создает структуру
func (g *JSONToStructGenerator) analyzeObject(obj map[string]interface{}, structName string) {
	if _, exists := g.structs[structName]; exists {
		return // Структура уже проанализирована
	}
	
	structInfo := &StructInfo{
		Name:   structName,
		Fields: make([]StructField, 0),
	}
	
	// Сортируем ключи для стабильного вывода
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	
	for _, key := range keys {
		value := obj[key]
		fieldName := toPascalCase(key)
		fieldType := g.analyzeValue(value, key)
		
		field := StructField{
			Name:     fieldName,
			Type:     fieldType,
			JSONTag:  key,
			Optional: value == nil,
		}
		
		structInfo.Fields = append(structInfo.Fields, field)
	}
	
	g.structs[structName] = structInfo
}

// GenerateFromJSON генерирует структуры из JSON данных
func (g *JSONToStructGenerator) GenerateFromJSON(data []byte, rootStructName string) error {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("не удалось распарсить JSON: %w", err)
	}
	
	switch v := jsonData.(type) {
	case map[string]interface{}:
		g.analyzeObject(v, rootStructName)
	case []interface{}:
		if len(v) > 0 {
			if obj, ok := v[0].(map[string]interface{}); ok {
				g.analyzeObject(obj, rootStructName+"Item")
			}
		}
	default:
		return fmt.Errorf("корневой элемент JSON должен быть объектом или массивом")
	}
	
	return nil
}

// GenerateGoCode генерирует Go код структур
func (g *JSONToStructGenerator) GenerateGoCode() string {
	var builder strings.Builder
	
	builder.WriteString("// Автоматически сгенерированные структуры из JSON\n\n")
	
	// Сортируем структуры по имени
	structNames := make([]string, 0, len(g.structs))
	for name := range g.structs {
		structNames = append(structNames, name)
	}
	sort.Strings(structNames)
	
	for _, name := range structNames {
		structInfo := g.structs[name]
		builder.WriteString(fmt.Sprintf("type %s struct {\n", structInfo.Name))
		
		for _, field := range structInfo.Fields {
			jsonTag := fmt.Sprintf("`json:\"%s", field.JSONTag)
			if field.Optional {
				jsonTag += ",omitempty"
			}
			jsonTag += "\"`"
			
			builder.WriteString(fmt.Sprintf("\t%s %s %s\n", 
				field.Name, field.Type, jsonTag))
		}
		
		builder.WriteString("}\n\n")
	}
	
	return builder.String()
}

func main() {
	filePath := "config.json"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	fmt.Printf("Анализ JSON файла: %s\n", filePath)
	
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		return
	}

	generator := NewJSONToStructGenerator()
	
	if err := generator.GenerateFromJSON(data, "Config"); err != nil {
		fmt.Printf("Ошибка анализа JSON: %v\n", err)
		return
	}

	fmt.Println("\nСгенерированные Go структуры:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Print(generator.GenerateGoCode())
	
	// Демонстрация использования
	fmt.Println("// Пример использования:")
	fmt.Println("func readConfig(filePath string) (*Config, error) {")
	fmt.Println("\tdata, err := os.ReadFile(filePath)")
	fmt.Println("\tif err != nil {")
	fmt.Println("\t\treturn nil, err")
	fmt.Println("\t}")
	fmt.Println("\t")
	fmt.Println("\tvar config Config")
	fmt.Println("\tif err := json.Unmarshal(data, &config); err != nil {")
	fmt.Println("\t\treturn nil, err")
	fmt.Println("\t}")
	fmt.Println("\t")
	fmt.Println("\treturn &config, nil")
	fmt.Println("}")
}
```

## 3. Универсальный конфигурационный ридер
```GO
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
```

## 4. Интеграция с основным файлом для работы со всеми типами конфигураций
```GO
// config_manager.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gcfg.v1"
)

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
	// Для демонстрации используем универсальную структуру
	config := make(map[string]map[string]interface{})
	
	// Читаем как обычный INI файл
	iniConfig := struct {
		Section map[string]map[string]string
	}{}
	
	if err := gcfg.ReadFileInto(&iniConfig, cm.FilePath); err != nil {
		return fmt.Errorf("ошибка чтения INI файла: %w", err)
	}
	
	cm.INIData = iniConfig
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
		if data, ok := cm.INIData.(struct{ Section map[string]map[string]string }); ok {
			for section, values := range data.Section {
				fmt.Printf("  [%s]\n", section)
				for key, value := range values {
					fmt.Printf("    %s = %s\n", key, value)
				}
			}
		}
	}
}

func main() {
	// Ищем все конфигурационные файлы в текущей директории
	configFiles, err := findConfigFiles(".")
	if err != nil {
		fmt.Printf("Ошибка поиска файлов: %v\n", err)
		return
	}
	
	if len(configFiles) == 0 {
		fmt.Println("Конфигурационные файлы не найдены")
		return
	}
	
	fmt.Printf("Найдено %d конфигурационных файлов:\n", len(configFiles))
	for i, file := range configFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}
	
	fmt.Println(strings.Repeat("=", 60))
	
	// Обрабатываем каждый файл
	manager := NewConfigManager()
	
	for _, filePath := range configFiles {
		fmt.Printf("\nОбработка файла: %s\n", filePath)
		fmt.Println(strings.Repeat("=", 60))
		
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

// findConfigFiles - функция из предыдущего примера
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
```

## Пример тестового JSON файла (test_config.json)
```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "myapp",
    "credentials": {
      "username": "admin",
      "password": "secret"
    },
    "ssl": true,
    "timeout": 30.5
  },
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "debug": true,
    "middlewares": ["cors", "auth", "logging"],
    "limits": {
      "max_connections": 1000,
      "request_timeout": 60,
      "body_size": "10MB"
    }
  },
  "logging": {
    "level": "info",
    "outputs": ["console", "file"],
    "file_config": {
      "path": "/var/log/app.log",
      "max_size": 100,
      "rotate": true
    }
  },
  "features": {
    "cache_enabled": true,
    "metrics_enabled": false,
    "experimental": ["feature_a", "feature_b"]
  },
  "version": "1.2.3",
  "environment": "development"
}
```

## 5. Расширенный генератор структур с поддержкой сложных случаев
```GO
// advanced_struct_generator.go 
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

// TypeInfo содержит информацию о типе
type TypeInfo struct {
	GoType      string
	IsOptional  bool
	IsArray     bool
	ElementType string
}

// AdvancedStructGenerator расширенный генератор структур
type AdvancedStructGenerator struct {
	structs     map[string]*StructInfo
	typeCounter map[string]int
}

// NewAdvancedStructGenerator создает новый расширенный генератор
func NewAdvancedStructGenerator() *AdvancedStructGenerator {
	return &AdvancedStructGenerator{
		structs:     make(map[string]*StructInfo),
		typeCounter: make(map[string]int),
	}
}

// analyzeTypeFromMultipleValues анализирует тип из нескольких значений (для массивов)
func (g *AdvancedStructGenerator) analyzeTypeFromMultipleValues(values []interface{}, fieldName string) TypeInfo {
	if len(values) == 0 {
		return TypeInfo{GoType: "interface{}", IsArray: true}
	}
	
	// Анализируем все элементы массива
	typeMap := make(map[string]int)
	var hasNil bool
	
	for _, value := range values {
		if value == nil {
			hasNil = true
			continue
		}
		
		typeInfo := g.analyzeValueAdvanced(value, fieldName+"Item")
		typeMap[typeInfo.GoType]++
	}
	
	// Если все элементы одного типа
	if len(typeMap) == 1 {
		for goType := range typeMap {
			return TypeInfo{
				GoType:      goType,
				IsArray:     true,
				ElementType: goType,
				IsOptional:  hasNil,
			}
		}
	}
	
	// Если разные типы, используем interface{}
	return TypeInfo{
		GoType:     "interface{}",
		IsArray:    true,
		IsOptional: hasNil,
	}
}

// analyzeValueAdvanced расширенный анализ значений
func (g *AdvancedStructGenerator) analyzeValueAdvanced(value interface{}, fieldName string) TypeInfo {
	switch v := value.(type) {
	case nil:
		return TypeInfo{GoType: "interface{}", IsOptional: true}
	case bool:
		return TypeInfo{GoType: "bool"}
	case float64:
		// Проверяем, является ли число целым
		if v == float64(int64(v)) && v >= -9223372036854775808 && v <= 9223372036854775807 {
			return TypeInfo{GoType: "int64"}
		}
		return TypeInfo{GoType: "float64"}
	case string:
		return TypeInfo{GoType: "string"}
	case []interface{}:
		if len(v) == 0 {
			return TypeInfo{GoType: "interface{}", IsArray: true}
		}
		return g.analyzeTypeFromMultipleValues(v, fieldName)
	case map[string]interface{}:
		structName := g.generateUniqueStructName(fieldName)
		g.analyzeObjectAdvanced(v, structName)
		return TypeInfo{GoType: structName}
	default:
		return TypeInfo{GoType: "interface{}"}
	}
}

// generateUniqueStructName генерирует уникальное имя структуры
func (g *AdvancedStructGenerator) generateUniqueStructName(baseName string) string {
	structName := toPascalCase(baseName)
	if structName == "" {
		structName = "NestedStruct"
	}
	
	// Проверяем, существует ли уже такое имя
	if _, exists := g.structs[structName]; !exists {
		return structName
	}
	
	// Генерируем уникальное имя
	g.typeCounter[structName]++
	return fmt.Sprintf("%s%d", structName, g.typeCounter[structName])
}

// analyzeObjectAdvanced расширенный анализ объектов
func (g *AdvancedStructGenerator) analyzeObjectAdvanced(obj map[string]interface{}, structName string) {
	if _, exists := g.structs[structName]; exists {
		return
	}
	
	structInfo := &StructInfo{
		Name:   structName,
		Fields: make([]StructField, 0),
	}
	
	// Сортируем ключи
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	
	for _, key := range keys {
		value := obj[key]
		fieldName := toPascalCase(key)
		typeInfo := g.analyzeValueAdvanced(value, key)
		
		goType := typeInfo.GoType
		if typeInfo.IsArray {
			goType = "[]" + typeInfo.GoType
		}
		
		// Добавляем указатель для опциональных полей
		if typeInfo.IsOptional && !typeInfo.IsArray {
			goType = "*" + goType
		}
		
		field := StructField{
			Name:     fieldName,
			Type:     goType,
			JSONTag:  key,
			Optional: typeInfo.IsOptional,
		}
		
		structInfo.Fields = append(structInfo.Fields, field)
	}
	
	g.structs[structName] = structInfo
}

// GenerateFromJSONAdvanced генерирует структуры с расширенным анализом
func (g *AdvancedStructGenerator) GenerateFromJSONAdvanced(data []byte, rootStructName string) error {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("не удалось распарсить JSON: %w", err)
	}
	
	switch v := jsonData.(type) {
	case map[string]interface{}:
		g.analyzeObjectAdvanced(v, rootStructName)
	case []interface{}:
		if len(v) > 0 {
			// Анализируем все элементы массива для более точного определения типа
			for i, item := range v {
				if obj, ok := item.(map[string]interface{}); ok {
					itemStructName := fmt.Sprintf("%sItem", rootStructName)
					if i > 0 {
						itemStructName = fmt.Sprintf("%sItem%d", rootStructName, i)
					}
					g.analyzeObjectAdvanced(obj, itemStructName)
				}
			}
		}
	default:
		return fmt.Errorf("корневой элемент JSON должен быть объектом или массивом")
	}
	
	return nil
}

// GenerateGoCodeAdvanced генерирует улучшенный Go код
func (g *AdvancedStructGenerator) GenerateGoCodeAdvanced(packageName string) string {
	var builder strings.Builder
	
	builder.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	builder.WriteString("// Автоматически сгенерированные структуры из JSON\n")
	builder.WriteString("// Сгенерировано с помощью AdvancedStructGenerator\n\n")
	
	// Добавляем импорты если нужны
	builder.WriteString("import (\n")
	builder.WriteString("\t\"encoding/json\"\n")
	builder.WriteString(")\n\n")
	
	// Сортируем структуры
	structNames := make([]string, 0, len(g.structs))
	for name := range g.structs {
		structNames = append(structNames, name)
	}
	sort.Strings(structNames)
	
	for _, name := range structNames {
		structInfo := g.structs[name]
		
		// Добавляем комментарий к структуре
		builder.WriteString(fmt.Sprintf("// %s представляет конфигурационные данные\n", structInfo.Name))
		builder.WriteString(fmt.Sprintf("type %s struct {\n", structInfo.Name))
		
		for _, field := range structInfo.Fields {
			// Добавляем комментарий к полю
			builder.WriteString(fmt.Sprintf("\t// %s соответствует JSON полю \"%s\"\n", field.Name, field.JSONTag))
			
			jsonTag := fmt.Sprintf("`json:\"%s", field.JSONTag)
			if field.Optional {
				jsonTag += ",omitempty"
			}
			jsonTag += "\"`"
			
			builder.WriteString(fmt.Sprintf("\t%s %s %s\n", 
				field.Name, field.Type, jsonTag))
		}
		
		builder.WriteString("}\n\n")
		
		// Генерируем методы для структуры
		g.generateMethods(&builder, structInfo)
	}
	
	return builder.String()
}

// generateMethods генерирует полезные методы для структур
func (g *AdvancedStructGenerator) generateMethods(builder *strings.Builder, structInfo *StructInfo) {
	structName := structInfo.Name
	
	// Метод для загрузки из JSON файла
	builder.WriteString(fmt.Sprintf("// Load%sFromFile загружает %s из JSON файла\n", structName, structName))
	builder.WriteString(fmt.Sprintf("func Load%sFromFile(filename string) (*%s, error) {\n", structName, structName))
	builder.WriteString("\tdata, err := os.ReadFile(filename)\n")
	builder.WriteString("\tif err != nil {\n")
	builder.WriteString("\t\treturn nil, err\n")
	builder.WriteString("\t}\n\n")
	builder.WriteString(fmt.Sprintf("\tvar config %s\n", structName))
	builder.WriteString("\tif err := json.Unmarshal(data, &config); err != nil {\n")
	builder.WriteString("\t\treturn nil, err\n")
	builder.WriteString("\t}\n\n")
	builder.WriteString("\treturn &config, nil\n")
	builder.WriteString("}\n\n")
	
	// Метод для сохранения в JSON файл
	builder.WriteString(fmt.Sprintf("// SaveToFile сохраняет %s в JSON файл\n", structName))
	builder.WriteString(fmt.Sprintf("func (c *%s) SaveToFile(filename string) error {\n", structName))
	builder.WriteString("\tdata, err := json.MarshalIndent(c, \"\", \"  \")\n")
	builder.WriteString("\tif err != nil {\n")
	builder.WriteString("\t\treturn err\n")
	builder.WriteString("\t}\n\n")
	builder.WriteString("\treturn os.WriteFile(filename, data, 0644)\n")
	builder.WriteString("}\n\n")
	
	// Метод для валидации (базовый)
	builder.WriteString(fmt.Sprintf("// Validate выполняет базовую валидацию %s\n", structName))
	builder.WriteString(fmt.Sprintf("func (c *%s) Validate() error {\n", structName))
	builder.WriteString("\t// TODO: Добавьте свою логику валидации\n")
	builder.WriteString("\treturn nil\n")
	builder.WriteString("}\n\n")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run advanced_struct_generator.go <json_file> [package_name] [output_file]")
		fmt.Println("Пример: go run advanced_struct_generator.go config.json config config.go")
		return
	}
	
	filePath := os.Args[1]
	packageName := "main"
	outputFile := ""
	
	if len(os.Args) > 2 {
		packageName = os.Args[2]
	}
	if len(os.Args) > 3 {
		outputFile = os.Args[3]
	}
	
	fmt.Printf("Анализ JSON файла: %s\n", filePath)
	
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		return
	}

	generator := NewAdvancedStructGenerator()
	
	if err := generator.GenerateFromJSONAdvanced(data, "Config"); err != nil {
		fmt.Printf("Ошибка анализа JSON: %v\n", err)
		return
	}

	goCode := generator.GenerateGoCodeAdvanced(packageName)
	
	if outputFile != "" {
		// Сохраняем в файл
		if err := os.WriteFile(outputFile, []byte(goCode), 0644); err != nil {
			fmt.Printf("Ошибка записи в файл %s: %v\n", outputFile, err)
			return
		}
		fmt.Printf("Go код сохранен в файл: %s\n", outputFile)
	} else {
		// Выводим в консоль
		fmt.Println("\nСгенерированный Go код:")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Print(goCode)
	}
	
	fmt.Println("Генерация завершена успешно!")
}
```

## 6. Пример использования всех компонентов
```GO
// example_usage.go 
package main

import (
	"fmt"
	"os"
)

func main() {
	// Создаем тестовый JSON файл
	testJSON := `{
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "myapp",
    "credentials": {
      "username": "admin",
      "password": "secret"
    },
    "ssl": true,
    "timeout": 30.5
  },
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "debug": true,
    "middlewares": ["cors", "auth", "logging"],
    "limits": {
      "max_connections": 1000,
      "request_timeout": 60,
      "body_size": "10MB"
    }
  },
  "logging": {
    "level": "info",
    "outputs": ["console", "file"],
    "file_config": {
      "path": "/var/log/app.log",
      "max_size": 100,
      "rotate": true
    }
  },
  "features": {
    "cache_enabled": true,
    "metrics_enabled": false,
    "experimental": ["feature_a", "feature_b"]
  },
  "version": "1.2.3",
  "environment": "development"
}`

	// Сохраняем тестовый файл
	if err := os.WriteFile("test_config.json", []byte(testJSON), 0644); err != nil {
		fmt.Printf("Ошибка создания тестового файла: %v\n", err)
		return
	}

	fmt.Println("Создан тестовый файл: test_config.json")
	
	// Демонстрируем работу с универсальным ридером
	fmt.Println("\n=== Демонстрация универсального ридера ===")
	reader := NewConfigReader()
	if err := reader.ReadJSON("test_config.json"); err != nil {
		fmt.Printf("Ошибка чтения: %v\n", err)
		return
	}

	// Получаем различные значения
	if host, err := reader.GetString("database.host"); err == nil {
		fmt.Printf("Database host: %s\n", host)
	}

	if port, err := reader.GetInt("database.port"); err == nil {
		fmt.Printf("Database port: %d\n", port)
	}

	if ssl, err := reader.GetBool("database.ssl"); err == nil {
		fmt.Printf("Database SSL: %t\n", ssl)
	}

	if middlewares, err := reader.GetStringArray("server.middlewares"); err == nil {
		fmt.Printf("Server middlewares: %v\n", middlewares)
	}

	if timeout, err := reader.GetFloat("database.timeout"); err == nil {
		fmt.Printf("Database timeout: %g\n", timeout)
	}

	// Получаем вложенный объект
	if creds, err := reader.GetObject("database.credentials"); err == nil {
		fmt.Printf("Database credentials: %v\n", creds)
	}

	fmt.Println("\nВсе доступные ключи:")
	keys := reader.GetAllKeys()
	for i, key := range keys {
		if value, exists := reader.Get(key); exists {
			fmt.Printf("%d. %s (%T)\n", i+1, key, value)
		}
	}
}
```

## Команды для запуска:
```bash
go run cmd/work-configs/dynamic_json_reader.go test_config.json

go run cmd/work-configs/json_to_struct.go test_config.json

go run cmd/work-configs/advanced_struct_generator.go test_config.json config generated_config.go

go run cmd/work-configs/universal_config_reader.go test_config.json

go run cmd/work-configs/config_manager.go

go run cmd/work-configs/example_usage.go
```

## Основные возможности решения:
* Динамическое чтение JSON - без предварительного знания структуры
* Автогенерация Go структур - создание типизированных структур из JSON
* Универсальный доступ к данным - получение значений по ключам с проверкой типов
* Поддержка вложенных объектов - работа с многоуровневыми конфигурациями
* Генерация методов - автоматическое создание Load/Save методов
* Обработка массивов - корректная работа с массивами разных типов
* Валидация типов - проверка соответствия типов при получении значений

    Теперь у вас есть полный набор инструментов для работы с JSON конфигурациями любой структуры!