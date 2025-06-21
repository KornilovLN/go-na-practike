package utils

// finder.go

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/KornilovLN/go-na-practike/cmd/wrk-configs/pkg/types"
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
