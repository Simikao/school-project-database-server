package main

import (
	"context"
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
	admins := redditos.Collection("admins")
	comments := redditos.Collection("comments")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) { handler.MainPage(c, posts) })

	r.GET("/u", func(c *gin.Context) { handler.GetUsers(c, users) })
	r.POST("/u", func(c *gin.Context) { handler.AddNewUser(c, users) })
	r.GET("/u/:name", func(c *gin.Context) { handler.GetUser(c, users) })
	r.PUT("/u/:name", func(c *gin.Context) { handler.UpdateUser(c, users) })
	r.POST("/a", func(c *gin.Context) { handler.AddNewAdmin(c, users, admins) })
	r.DELETE("/u", func(c *gin.Context) { handler.DeleteUser(c, users) })
	r.GET("/u/search/:name", func(c *gin.Context) { handler.SearchUsers(c, users) })
	r.GET("/u/top", func(c *gin.Context) { handler.ShowBestUsers(c, users, posts, comments) })

	r.POST("/new-post", func(c *gin.Context) { handler.AddNewPost(c, posts) })

	r.POST("/new-community", func(c *gin.Context) { handler.AddNewCommunity(c, communities) })
	r.GET("/c", func(c *gin.Context) { handler.GetCommunities(c, communities) })

	initializers.OGAdmin(admins)
	go r.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	quitCtx, quitCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer quitCancel()
	// err = users.Drop(quitCtx)
	// fmt.Println()
	// log.Info("Dropping users")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// err = posts.Drop(quitCtx)
	// fmt.Println()
	// log.Info("Dropping posts")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// err = communities.Drop(quitCtx)
	// fmt.Println()
	// log.Info("Dropping communities")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	log.Info("Closing connection")
	err = users.Database().Client().Disconnect(quitCtx)
	if err != nil {
		log.Fatal(err.Error())
	}

}
