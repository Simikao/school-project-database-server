package handler

import (
	"context"
	"encoding/json"
	"io"
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
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("jsonDataError:" + err.Error())
	}
	log.Info(string(jsonData))

	err = json.Unmarshal(jsonData, &user)
	log.Debug(user)
	if err != nil {
		log.Error("Error unmarshalling:" + err.Error())
		c.JSON(400, struct {
			Err string `json:"error"`
			Msg string `json:"message"`
		}{"Invalid syntax", "Cannot read the JSON file"})
		return
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
