package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Speed    string `json:"speed"`
	IP       string `json:"ip"`
}

func LoadConfig(path string) (*Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("No se puede leer el archivo %w", err)
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Formato JSON invalido: %w", err)
	}

	if config.IP == "" || config.User == "" {
		return nil, fmt.Errorf("El archivo contiene campos obligatorios vacios")
	}

	return &config, nil
}
