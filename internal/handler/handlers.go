package handler

import (
	"context"
	"encoding/json"
	"net/http"
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
		c.JSON(500, datatype.Response{
			Success: false,
			Data:    "Something went wrong with age validator",
		})
	}

	err = validate.RegisterValidation("isUnique", validators.IsUniqueName)
	if err != nil {
		c.JSON(500, datatype.Response{
			Success: false,
			Data:    "Something went wrong with name validator",
		})
	}

	var user datatype.User
	err = json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Badly formatted JSON",
		})
	}

	err = validate.Struct(user)
	if err != nil {
		log.Error(err.Error())
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    errors.Error(),
		})
		return
	}

	id, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatype.Response{
			Success: false,
			Data:    "Problem with database",
		})
		return
	}

	response := id.InsertedID.(primitive.ObjectID).Hex()
	log.Debug("Added user of id:" + response)
	c.JSON(http.StatusCreated, datatype.Response{
		Success: true,
		Data:    "User added, welcome " + user.Name,
	})
}

func AddNewPost(c *gin.Context, collection *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	validate := validator.New()

	err := validate.RegisterValidation("isUnique", validators.IsUniqueTitle)
	if err != nil {
		c.JSON(500, datatype.Response{
			Success: false,
			Data:    "Something went wrong with name validator",
		})
	}

	var post datatype.Post
	err = json.NewDecoder(c.Request.Body).Decode(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Badly formatted JSON",
		})
	}

	err = validate.Struct(post)
	if err != nil {
		log.Error(err.Error())
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    errors.Error(),
		})
		return
	}

	id, err := collection.InsertOne(ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatype.Response{
			Success: false,
			Data:    "Problem with database",
		})
		return
	}

	response := id.InsertedID.(primitive.ObjectID).Hex()
	log.Debug("Added post of id:" + response)
	c.JSON(http.StatusCreated, datatype.Response{
		Success: true,
		Data:    "Post added",
	})
}

func AddNewCommunity(c *gin.Context, collection *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	validate := validator.New()
	err := validate.RegisterValidation("isUnique", validators.IsUniqueName)
	if err != nil {
		log.Error(err.Error())
	}
	var community datatype.Community
	err = json.NewDecoder(c.Request.Body).Decode(&community)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Badly formatted JSON",
		})
	}

	err = validate.Struct(community)
	if err != nil {
		log.Error(err.Error())
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    errors.Error(),
		})
		return
	}

	id, err := collection.InsertOne(ctx, community)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatype.Response{
			Success: false,
			Data:    "Problem with database",
		})
		return
	}

	response := id.InsertedID.(primitive.ObjectID).Hex()
	log.Debug("Added post of id:" + response)
	c.JSON(http.StatusCreated, datatype.Response{
		Success: true,
		Data:    "Added community " + community.Name,
	})
}
