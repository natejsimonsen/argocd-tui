package config

import (
	"fmt"
	"log"
	"os"

	"example.com/main/services/utils"
	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v3"
)

type InternalConfig struct {
	Foreground string `yaml:"foreground"`
	Background string `yaml:"background"`
}

type Config struct {
	Background tcell.Color
}

func NewConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error loading os home dir: %v", err)
	}
	path := fmt.Sprintf("%s/.config/argocd-tui/config.yaml", home)

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var config InternalConfig

	err = yaml.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling yaml: %v", err)
	}

	externalConfig := Config{
		Background: utils.HexToColor(config.Background, tcell.ColorSkyblue),
	}

	return &externalConfig
}
