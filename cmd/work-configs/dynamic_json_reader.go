// dynamic_json_reader.go
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
	filePath := "conf.json"
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
