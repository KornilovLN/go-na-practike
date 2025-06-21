// json_to_struct.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	filePath := "cmd/wrk-configs/configs/examples/conf.json"
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
