package parsers

// yaml.go

import (
	"os"

	"github.com/KornilovLN/go-na-praktike/cmd/wrk-configs/pkg/types"
	"gopkg.in/yaml.v3"
)

type YAMLParser struct{}

func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

func (p *YAMLParser) Parse(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (p *YAMLParser) ParseFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return p.Parse(data, v)
}

func (p *YAMLParser) Format() types.ConfigFormat {
	return types.FormatYAML
}
