package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	DSToken      string `json:"ds_token"`
	SettingsFile string `json:"settings_file"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, errors.New("не удалось открыть файл конфигурации: " + err.Error())
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, errors.New("ошибка разбора JSON: " + err.Error())
	}

	if config.DSToken == "" {
		return nil, errors.New("не указан токен Discord бота")
	}
	return &config, nil
}
