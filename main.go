package main

import (
	"context"
	"fmt"
	"net/mail"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Name     string `json:"name" bson:"name"`
}

type Post struct {
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
	Author  User   `json:"author" bson:"author"`
}

var (
	uri  = "mongodb://localhost:27017"
	auth = options.Credential{
		Username: "root",
		Password: "example",
	}
)

func main() {
	log.Info("Hello world")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(auth))
	if err != nil {
		log.Fatal("Set up failed: ", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Connection failed: ", err)
	}

	redditos := client.Database("redditos")
	ex := redditos.Collection("example")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello World")
	})
	r.POST("/", func(c *gin.Context) {
		u := User{
			Email:    "ala@niemako.ta",
			Name:     "Adam",
			Password: "hashysz",
		}
		post := Post{
			Title:   "Bigsmall World",
			Content: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			Author:  u,
		}
		c.JSON(200, post)
	})
	r.POST("/add", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		newUser := User{
			Email:    "ala@niemako.ta",
			Name:     "Adam",
			Password: "hashysz",
		}
		_, err := mail.ParseAddress(newUser.Email)
		if err != nil {
			c.String(200, "Wrong email address")
		}
		id, err := ex.InsertOne(ctx, newUser)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.String(200, id.InsertedID.(primitive.ObjectID).Hex())
		var result User
		err = ex.FindOne(ctx, bson.M{"_id": id.InsertedID}).Decode(&result)
		if err != nil {
			c.String(500, "Something went wrong")
			return
		}
		c.JSON(200, result)

	})
	r.GET("/find", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		cur, err := ex.Find(ctx, bson.D{})
		if err != nil {
			c.String(500, err.Error())
		}
		var results []User
		for cur.Next(ctx) {
			var elem User
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
	err = ex.Drop(quitCtx)
	fmt.Println("Dropping database")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Closing connection")
	err = ex.Database().Client().Disconnect(quitCtx)
	if err != nil {
		log.Fatal(err.Error())
	}

}
