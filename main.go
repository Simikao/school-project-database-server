package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/datatype"
	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/handler"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri  = "mongodb://localhost:27017"
	auth = options.Credential{
		Username: "root",
		Password: "example",
	}
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

func main() {
	logLevel()

	// newUser := datatype.User{
	// 	Email:    "ala@niemako.ta",
	// 	Name:     "Adamek",
	// 	Password: "hashysz",
	// 	DoB:      time.Date(2020, time.Now().Month(), 2, 0, 0, 0, 0, time.Local),
	// }
	log.SetStyles(newStyle())
	log.Info("Hello world")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(auth))
	if err != nil {
		log.Fatal("Set up failed: ", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Connection failed: ", "details -> ", err)
	}

	redditos := client.Database("redditos")
	users := redditos.Collection("users")
	posts := redditos.Collection("posts")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello World")
	})

	r.POST("/add", func(c *gin.Context) { handler.AddNewUser(c, users) })
	r.POST("/new-post", func(c *gin.Context) { handler.AddNewPost(c, posts) })

	r.GET("/find", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		cur, err := users.Find(ctx, bson.D{})
		if err != nil {
			c.String(500, err.Error())
		}
		var results []datatype.User
		for cur.Next(ctx) {
			var elem datatype.User
			err := cur.Decode(&elem)
			if err != nil {
				c.String(500, err.Error())
			}
			results = append(results, elem)

		}
		c.JSON(200, results)

	})

	go r.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	quitCtx, quitCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer quitCancel()
	err = users.Drop(quitCtx)
	fmt.Println()
	log.Info("Dropping database")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Closing connection")
	err = users.Database().Client().Disconnect(quitCtx)
	if err != nil {
		log.Fatal(err.Error())
	}

}
