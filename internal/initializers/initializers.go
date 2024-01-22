package initializers

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func newStyle() (style *log.Styles) {
	style = log.DefaultStyles()
	pinkText := lipgloss.NewStyle().Foreground(lipgloss.Color("#c88d94"))
	greyText := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	style.Key = pinkText
	style.Value = greyText
	return
}

func logLevel() {
	switch level := os.Getenv("SERVERLOG"); level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func PreStart() {
	logLevel()

	log.SetStyles(newStyle())

	log.Info("Hello world")
}
