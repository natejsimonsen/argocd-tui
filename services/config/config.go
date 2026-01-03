package config

import (
	"fmt"
	"log"
	"os"

	"example.com/main/services/utils"
	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v3"
)

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
		Background:  utils.HexToColor(config.Colors.Background, tcell.ColorBlack),
		Border:      utils.HexToColor(config.Colors.Border, tcell.ColorDarkSlateGray),
		Selected:    utils.HexToColor(config.Colors.Selected, tcell.ColorSkyblue),
		Header:      utils.HexToColor(config.Colors.Header, tcell.ColorGray),
		Text:        utils.HexToColor(config.Colors.Text, tcell.ColorWhite),
		Foreground:  utils.HexToColor(config.Colors.Foreground, tcell.ColorWhiteSmoke),
		Progressing: utils.HexToColor(config.Colors.Progressing, tcell.ColorLightBlue),
		Missing:     utils.HexToColor(config.Colors.Missing, tcell.ColorLightYellow),
		Healthy:     utils.HexToColor(config.Colors.Healthy, tcell.ColorLightGreen),
		Degraded:    utils.HexToColor(config.Colors.Degraded, tcell.ColorRed),
	}

	return &externalConfig
}
