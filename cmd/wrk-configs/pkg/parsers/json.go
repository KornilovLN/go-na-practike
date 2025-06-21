package parsers

// json.go

import (
	"encoding/json"
	"os"

	"github.com/KornilovLN/go-na-practike/cmd/wrk-configs/pkg/types"
)

type JSONParser struct{}

func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

func (p *JSONParser) Parse(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (p *JSONParser) ParseFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return p.Parse(data, v)
}

func (p *JSONParser) Format() types.ConfigFormat {
	return types.FormatJSON
}

// ParseDynamic парсит JSON в map[string]interface{} для динамического доступа
func (p *JSONParser) ParseDynamic(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	return result, err
}

// ParseDynamicFile парсит JSON файл в map[string]interface{}
func (p *JSONParser) ParseDynamicFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return p.ParseDynamic(data)
}
