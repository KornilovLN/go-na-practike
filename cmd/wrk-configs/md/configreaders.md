# Предлагаемая структура проекта:
```
cmd/work-configs/
├── README.md                    # Описание всех примеров
├── configs/                     # Все тестовые конфигурационные файлы
│   ├── examples/
│   │   ├── conf.ini
│   │   ├── conf.json
│   │   ├── conf.yaml
│   │   ├── test_config.json
│   │   └── test2_config.json
│   └── schemas/                 # JSON схемы для валидации
├── pkg/                         # Переиспользуемые пакеты
│   ├── parsers/                 # Парсеры для разных форматов
│   │   ├── json.go
│   │   ├── yaml.go
│   │   ├── ini.go
│   │   └── universal.go
│   ├── generators/              # Генераторы конфигов
│   │   ├── json.go
│   │   ├── yaml.go
│   │   └── ini.go
│   ├── utils/                   # Утилиты
│   │   ├── finder.go           # Поиск конфиг файлов
│   │   └── validator.go        # Валидация
│   └── types/                   # Общие типы и структуры
│       └── config.go
├── examples/                    # Примеры использования
│   ├── 01-basic-json/
│   │   └── main.go
│   ├── 02-basic-yaml/
│   │   └── main.go
│   ├── 03-basic-ini/
│   │   └── main.go
│   ├── 04-dynamic-json/
│   │   └── main.go
│   ├── 05-universal-reader/
│   │   └── main.go
│   ├── 06-config-manager/
│   │   └── main.go
│   └── 07-json-to-struct/
│       └── main.go
├── cmd/                         # Исполняемые утилиты
│   ├── config-converter/
│   │   └── main.go
│   └── config-validator/
│       └── main.go
└── internal/                    # Внутренние пакеты
    └── author/
        └── author.go
```

## Команды для реорганизации:
```bash
mkdir -p cmd/work-configs/{configs/{examples,schemas},pkg/{parsers,generators,utils,types},examples/{01-basic-json,02-basic-yaml,03-basic-ini,04-dynamic-json,05-universal-reader,06-config-manager,07-json-to-struct},cmd/{config-converter,config-validator},internal/author}

cd cmd/work-configs

mv conf.* test*.json configs/examples/

mv author internal/
```

## Создание базовых файлов:
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
```

```Go
// finder.go
package utils

import (
	"os"
	"path/filepath"
	"strings"
	"time"

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
		
		format := getFormatByExtension(filepath.Ext(path))
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

func getFormatByExtension(ext string) types.ConfigFormat {
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
```

```GO
// main.go
package main

import (
	"fmt"
	"log"

	"github.com/KornilovLN/go-na-praktike/cmd/work-configs/pkg/parsers"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
	} `json:"database"`
	Debug bool `json:"debug"`
}

func main() {
	fmt.Println("=== Пример 1: Базовая работа с JSON ===")
	
	parser := parsers.NewJSONParser()
	
	var config Config
	err := parser.ParseFile("../../configs/examples/conf.json", &config)
	if err != nil {
		log.Fatal("Ошибка парсинга:", err)
	}
	
	fmt.Printf("Конфигурация загружена: %+v\n", config)
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
```

# Работа с конфигурационными файлами

    Этот проект содержит примеры различных подходов к парсингу и генерации конфигурационных файлов в Go.

## Структура проекта

- `pkg/` - переиспользуемые пакеты
  - `parsers/` - парсеры для разных форматов (JSON, YAML, INI)
  - `generators/` - генераторы конфигураций
  - `utils/` - утилиты (поиск файлов, валидация)
  - `types/` - общие типы и интерфейсы
- `examples/` - примеры использования
- `configs/` - тестовые конфигурационные файлы
- `cmd/` - исполняемые утилиты
- `internal/` - внутренние пакеты

## Примеры

1. **01-basic-json** - базовая работа с JSON
2. **02-basic-yaml** - базовая работа с YAML
3. **03-basic-ini** - базовая работа с INI
4. **04-dynamic-json** - динамическое чтение JSON
5. **05-universal-reader** - универсальный читатель конфигов
6. **06-config-manager** - менеджер конфигураций
7. **07-json-to-struct** - генерация структур из JSON

## Запуск примеров

```bash
cd examples/01-basic-json && go run main.go
cd examples/02-basic-yaml && go run main.go
# и т.д.
```

## Утилиты

- `config-converter` - конвертация между форматами
- `config-validator` - валидация конфигураций


## README.md
* **Теперь:**
  * Каждый пример изолирован в своей папке
  * Переиспользуемый код вынесен в pkg/
  * Тестовые конфиги собраны в одном месте
  * Легко добавлять новые примеры без затрагивания существующих
  * Четкое разделение между библиотечным кодом и примерами
  