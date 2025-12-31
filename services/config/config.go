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

const ARGO_CONFIG_DIR = "argocd-tui"

func NewConfig() *Config {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Could not find config dir: %v", err)
	}

	path := fmt.Sprintf("%s/argocd-tui/config.yaml", configDir)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(fmt.Sprintf("%s/%s", configDir, ARGO_CONFIG_DIR), 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}

		file, err := os.Create(path)
		if err != nil {
			log.Fatalf("Error creating file %s: %v", path, err)
		}

		defer file.Close()
	}

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
