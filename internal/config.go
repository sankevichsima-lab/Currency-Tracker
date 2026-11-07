package internal

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ServerConfig stores HTTP server network settings.
// ServerConfig хранит сетевые настройки HTTP-сервера.
type ServerConfig struct {
	Port string `yaml:"port"`
}

// DatabaseConfig stores database connection settings.
// DatabaseConfig хранит параметры подключения к базе данных.
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// APIConfig stores third-party API settings and polling interval.
// APIConfig хранит настройки стороннего API и интервал опроса.
type APIConfig struct {
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
}

// Config aggregates all application configuration sections.
// Config объединяет все секции конфигурации приложения.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	API      APIConfig      `yaml:"api"`
	Database DatabaseConfig `yaml:"database"`
}

// LoadConfig reads and parses the YAML configuration file.
// LoadConfig считывает и парсит файл конфигурации в формате YAML.
func LoadConfig(filepath string) (*Config, error) {
	// Read raw config file content from disk.
	// Считываем содержимое файла конфигурации с диска.
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Unmarshal YAML content into the Config struct.
	// Десериализуем YAML в структуру Config.
	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("deserialization issues: %w", err)
	}

	return &cfg, nil
}
