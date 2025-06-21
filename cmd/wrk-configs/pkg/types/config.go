// config.go
package types

import "time"

// ConfigFormat представляет тип конфигурационного файла
type ConfigFormat string

const (
	FormatJSON ConfigFormat = "json"
	FormatYAML ConfigFormat = "yaml"
	FormatINI  ConfigFormat = "ini"
	FormatTOML ConfigFormat = "toml"
)

// ConfigFile представляет информацию о конфигурационном файле
type ConfigFile struct {
	Path     string
	Format   ConfigFormat
	Size     int64
	Modified time.Time
}

// Parser интерфейс для парсеров конфигураций
type Parser interface {
	Parse(data []byte, v interface{}) error
	Format() ConfigFormat
}

// Generator интерфейс для генераторов конфигураций
type Generator interface {
	Generate(v interface{}) ([]byte, error)
	Format() ConfigFormat
}

// ConfigManager интерфейс для управления конфигурациями
type ConfigManager interface {
	Load(path string, v interface{}) error
	Save(path string, v interface{}) error
	Validate(path string) error
	Convert(srcPath, dstPath string, dstFormat ConfigFormat) error
}

// CommonConfig базовая структура конфигурации для примеров
type CommonConfig struct {
	Database struct {
		Host     string `json:"host" yaml:"host" ini:"host"`
		Port     int    `json:"port" yaml:"port" ini:"port"`
		Username string `json:"username" yaml:"username" ini:"username"`
		Password string `json:"password" yaml:"password" ini:"password"`
	} `json:"database" yaml:"database" ini:"database"`

	Server struct {
		Host string `json:"host" yaml:"host" ini:"host"`
		Port int    `json:"port" yaml:"port" ini:"port"`
	} `json:"server" yaml:"server" ini:"server"`

	Debug   bool `json:"debug" yaml:"debug" ini:"debug"`
	Logging struct {
		Level string `json:"level" yaml:"level" ini:"level"`
		File  string `json:"file" yaml:"file" ini:"file"`
	} `json:"logging" yaml:"logging" ini:"logging"`
}
