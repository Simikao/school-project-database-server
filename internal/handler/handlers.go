package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/datatype"
	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/validators"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ()

func AddNewUser(c *gin.Context, collection *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	validate := validator.New()
	err := validate.RegisterValidation("dob", validators.DOBValidator)
	if err != nil {
		log.Error(err.Error())
	}

	var user datatype.User
	err = json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		log.Error(err.Error())
	}

	err = validate.Struct(user)
	if err != nil {
		log.Error(err.Error())
		errors := err.(validator.ValidationErrors)
		c.String(400, errors.Error())
		return
	}

	id, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	response := id.InsertedID.(primitive.ObjectID).Hex()
	c.JSON(200, struct {
		Success bool   `json:"success"`
		Data    string `json:"payload"`
	}{
		true,
		"Added user of id: " + response,
	})
}