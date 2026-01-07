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

	config := &Config{}

	if len(os.Args) > 1 {
		path = os.Args[1]
		file, err := os.ReadFile(path)
		if err == nil {
			json.Unmarshal(file, config)
			fmt.Println("Configuration cargada con exito!", path)
			return config, nil
		}
		fmt.Println("No se pudo leer el archivo de configuración, usando valores por defecto.")
	}

	return config, nil
}
