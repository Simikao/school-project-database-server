package initializers

import (
	"context"
	"os"
	"time"

	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/datatype"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var Admins []string

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
	log.Debug("Running in debug mode")
}

func OGAdmin(admins *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	Admins = append(Admins, "Smith")

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 4)
	if err != nil {
		log.Error(err)
		return
	}
	strID := "6f0d4e28cb46a7a22c1a5e14"
	objID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		log.Error(err)
		return
	}

	newAdmin := datatype.Administrator{
		UserID:   objID,
		Name:     "Smith",
		Password: string(hash),
	}

	_, err = admins.InsertOne(ctx, newAdmin)
	if err != nil {
		log.Error(err, "Problem with database")
		return
	}

	log.Debug("Admin added")
}
