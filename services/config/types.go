package config

import (
	"github.com/gdamore/tcell/v2"
)

type InternalConfig struct {
	Colors struct {
		Text        string `yaml:"text"`
		Border      string `yaml:"text"`
		Header      string `yaml:"header"`
		Foreground  string `yaml:"foreground"`
		Selected    string `yaml:"selected"`
		Background  string `yaml:"background"`
		Progressing string `yaml:"progressing"`
		Missing     string `yaml:"missing"`
		Healthy     string `yaml:"healthy"`
		Degraded    string `yaml:"degraded"`
	} `yaml:"colors"`
}

type Config struct {
	Background  tcell.Color
	Text        tcell.Color
	Border      tcell.Color
	Header      tcell.Color
	Selected    tcell.Color
	Foreground  tcell.Color
	Progressing tcell.Color
	Missing     tcell.Color
	Healthy     tcell.Color
	Degraded    tcell.Color
}
