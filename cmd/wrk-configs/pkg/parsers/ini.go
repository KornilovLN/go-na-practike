// ini.go
package parsers

import (
	"github.com/KornilovLN/go-na-practike/cmd/wrk-configs/pkg/types"
	"gopkg.in/ini.v1"
)

type INIParser struct{}

func NewINIParser() *INIParser {
	return &INIParser{}
}

func (p *INIParser) Parse(data []byte, v interface{}) error {
	cfg, err := ini.Load(data)
	if err != nil {
		return err
	}
	return cfg.MapTo(v)
}

func (p *INIParser) ParseFile(path string, v interface{}) error {
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	return cfg.MapTo(v)
}

func (p *INIParser) Format() types.ConfigFormat {
	return types.FormatINI
}
