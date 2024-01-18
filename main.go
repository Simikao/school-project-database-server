package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
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
