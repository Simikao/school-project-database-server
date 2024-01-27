package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/datatype"
	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/handler"
	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/initializers"
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

func main() {

	initializers.PreStart()
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
	communities := redditos.Collection("communities")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusCreated, "hello World")
	})

	r.GET("/u", func(c *gin.Context) { handler.GetUsers(c, users) })
	r.POST("/u", func(c *gin.Context) { handler.AddNewUser(c, users) })
	r.GET("/u/:name", func(c *gin.Context) { handler.GetUser(c, users) })
	r.PUT("/u/:name", func(c *gin.Context) { handler.UpdateUser(c, users) })

	r.POST("/new-post", func(c *gin.Context) { handler.AddNewPost(c, posts) })

	r.POST("/new-community", func(c *gin.Context) { handler.AddNewCommunity(c, communities) })

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
	log.Info("Dropping users")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = posts.Drop(quitCtx)
	fmt.Println()
	log.Info("Dropping posts")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = communities.Drop(quitCtx)
	fmt.Println()
	log.Info("Dropping communities")
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Closing connection")
	err = users.Database().Client().Disconnect(quitCtx)
	if err != nil {
		log.Fatal(err.Error())
	}

}
