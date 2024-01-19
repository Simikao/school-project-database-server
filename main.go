package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	Uuid     uuid.UUID `json:"id"`
}

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  User   `json:"author"`
}

var uri = "mongodb://root:example@localhost:27017/?timeoutMS=1000"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	redditos := client.Database("redditos")
	ex := redditos.Collection("example")

	id, err := ex.InsertOne(ctx, bson.M{"foo": "bar", "dorrito": "dust"})
	if err != nil {
		log.Fatal(err)
	}

	result := ex.FindOne(ctx, bson.M{"_id": id.InsertedID})
	if err := result.Err(); err != nil {
		log.Fatal("50", err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		r := struct {
			Name  string `bson:"foo" json:"fiddle"`
			Name2 string `bson:"dorrito" json:"crisps"`
		}{}
		err := result.Decode(&r)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		c.JSON(200, r)
	})
	r.POST("/", func(c *gin.Context) {
		u := User{
			Email:    "ala@niemako.ta",
			Name:     "Adam",
			Password: "hashysz",
			Uuid:     uuid.New(),
		}
		post := Post{
			Title:   "Bigsmall World",
			Content: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			Author:  u,
		}
		c.JSON(200, post)
	})
	r.Run()
}
