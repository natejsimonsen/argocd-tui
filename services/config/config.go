package config

import (
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
	// TODO: make this load from home dir
	path := "/home/nate/.config/argocd-tui/config.yaml"

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
