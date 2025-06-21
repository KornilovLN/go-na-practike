# Создание новой структуры:
```bash
mkdir -p cmd/work-configs
cd cmd/work-configs
mkdir -p {configs/{examples,schemas},pkg/{parsers,generators,utils,types},examples/{01-basic-json,02-basic-yaml,03-basic-ini,04-dynamic-json,05-universal-reader,06-config-manager,07-json-to-struct},cmd/{config-converter,config-validator},internal}
```

## Создание основных файлов:
```GO
// go.mod
module github.com/KornilovLN/go-na-praktike/cmd/work-configs

go 1.21

require (
	gopkg.in/yaml.v3 v3.0.1
	gopkg.in/ini.v1 v1.67.0
)
```

```GO
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
	
	Debug   bool     `json:"debug" yaml:"debug" ini:"debug"`
	Logging struct {
		Level string `json:"level" yaml:"level" ini:"level"`
		File  string `json:"file" yaml:"file" ini:"file"`
	} `json:"logging" yaml:"logging" ini:"logging"`
}
```

```GO
// json.go 
package parsers

import (
	"encoding/json"
	"os"

	"github.com/KornilovLN/go-na-praktike/cmd/work-configs/pkg/types"
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
```

```GO
// yaml.go
package parsers

import (
	"os"

	"github.com/KornilovLN/go-na-praktike/cmd/work-configs/pkg/types"
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
```

```GO
// ini.go 
package parsers

import (
	"github.com/KornilovLN/go-na-praktike/cmd/work-configs/pkg/types"
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
```

```GO
// finder.go
package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/KornilovLN/go-na-praktike/cmd/work-configs/pkg/types"
)

// FindConfigFiles ищет конфигурационные файлы в указанной директории
func FindConfigFiles(dir string) ([]types.ConfigFile, error) {
	var configs []types.ConfigFile
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		format := GetFormatByExtension(filepath.Ext(path))
		if format != "" {
			configs = append(configs, types.ConfigFile{
				Path:     path,
				Format:   format,
				Size:     info.Size(),
				Modified: info.ModTime(),
			})
		}
		
		return nil
	})
	
	return configs, err
}

// GetFormatByExtension определяет формат по расширению файла
func GetFormatByExtension(ext string) types.ConfigFormat {
	switch strings.ToLower(ext) {
	case ".json":
		return types.FormatJSON
	case ".yaml", ".yml":
		return types.FormatYAML
	case ".ini":
		return types.FormatINI
	case ".toml":
		return types.FormatTOML
	default:
		return ""
	}
}

// GetFormatByContent пытается определить формат по содержимому файла
func GetFormatByContent(data []byte) types.ConfigFormat {
	content := strings.TrimSpace(string(data))
	
	if len(content) == 0 {
		return ""
	}
	
	// JSON начинается с { или [
	if content[0] == '{' || content[0] == '[' {
		return types.FormatJSON
	}
	
	// INI часто содержит секции [section]
	if strings.Contains(content, "[") && strings.Contains(content, "]") {
		return types.FormatINI
	}
	
	// YAML часто содержит : без кавычек
	if strings.Contains(content, ":") && !strings.Contains(content, `":`) {
		return types.FormatYAML
	}
	
	return ""
}
```

* **app.json**
```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "admin",
    "password": "secret123"
  },
  "server": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "debug": true,
  "logging": {
    "level": "info",
    "file": "/var/log/app.log"
  }
}
```

* **app.yaml**
```yml
database:
  host: localhost
  port: 5432
  username: admin
  password: secret123

server:
  host: 0.0.0.0
  port: 8080

debug: true

logging:
  level: info
  file: /var/log/app.log
```

* **app.ini**
```ini
[database]
host = localhost
port = 5432
username = admin
password = secret123

[server]
host = 0.0.0.0
port = 8080

debug = true

[logging]
level = info
file = /var/log/app.log
```

```GO
// main.go 
package main

import (
	"fmt"
	"log"

	"github.com/KornilovLN/go-na-praktike/cmd/work-configs/pkg/parsers"
	"github.com/KornilovLN/go-na-praktike/cmd/work-configs/pkg/types"
)

func main() {
	fmt.Println("=== Пример 1: Базовая работа с JSON ===")
	
	parser := parsers.NewJSONParser()
	
	var config types.CommonConfig
	err := parser.ParseFile("../../configs/examples/app.json", &config)
	if err != nil {
		log.Fatal("Ошибка парсинга:", err)
	}
	
	fmt.Printf("Конфигурация загружена:\n")
	fmt.Printf("  Database: %s:%d (user: %s)\n", 
		config.Database.Host, config.Database.Port, config.Database.Username)
	fmt.Printf("  Server: %s:%d\n", 
		config.Server.Host, config.Server.Port)
	fmt.Printf("  Debug: %v\n", config.Debug)
	fmt.Printf("  Logging: %s -> %s\n", 
		config.Logging.Level, config.Logging.File)
}
```

# Работа с конфигурационными файлами

    Этот проект содержит примеры различных подходов к парсингу и генерации конфигурационных файлов в Go.

## Структура проекта

```
cmd/work-configs/
├── pkg/                         # Переиспользуемые пакеты
│   ├── parsers/                 # Парсеры для разных форматов
│   ├── generators/              # Генераторы конфигураций  
│   ├── utils/                   # Утилиты
│   └── types/                   # Общие типы и интерфейсы
├── examples/                    # Примеры использования
├── configs/                     # Тестовые конфигурационные файлы
├── cmd/                         # Исполняемые утилиты
└── internal/                    # Внутренние пакеты
```

## Примеры

1. **01-basic-json** - базовая работа с JSON
2. **02-basic-yaml** - базовая работа с YAML  
3. **03-basic-ini** - базовая работа с INI
4. **04-dynamic-json** - динамическое чтение JSON
5. **05-universal-reader** - универсальный читатель конфигов
6. **06-config-manager** - менеджер конфигураций
7. **07-json-to-struct** - генерация структур из JSON

## Быстрый старт

```bash
# Инициализация модуля
go mod tidy

# Запуск примера
cd examples/01-basic-json && go run main.go
```

## Миграция из старой структуры

    Старые наработки остаются в корне,
    новые разработки ведутся в этой структуре.
    Постепенно переносить и рефакторить код из старых файлов.

## Добавление нового примера

1. Создайте папку `examples/XX-new-example/`
2. Добавьте `main.go` с вашим кодом
3. При необходимости расширьте пакеты в `pkg/`
```

```bash
go mod tidy
```

README.md
```
Теперь есть чистая новая структура в cmd/wrk-configs/, а старые наработки остались нетронутыми. 

Можно:
    Постепенно переносить код из старых файлов
    Тестировать новые идеи в изолированных примерах
    Развивать переиспользуемые модули в pkg/
    Не ломать существующие наработки
```    