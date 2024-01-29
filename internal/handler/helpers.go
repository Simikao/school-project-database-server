package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/datatype"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func bodyDecoder(c *gin.Context, body interface{}) error {
	err := json.NewDecoder(c.Request.Body).Decode(body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Badly formated JSON",
		})
		return err
	}
	return nil

}

func findUserByName(c *gin.Context, collection *mongo.Collection, name string, val *datatype.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	filter := struct {
		Name string `bson:"name"`
	}{
		Name: name,
	}
	log.Debug(filter)
	result := collection.FindOne(ctx, filter)
	err := result.Decode(val)
	if err == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, datatype.Response{
			Success: false,
			Data:    "User not found",
		})
		return false
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, datatype.Response{
			Success: false,
			Data:    "Failed decoding result",
		})
		return false
	}
	return true
}

func findUserByID(c *gin.Context, collection *mongo.Collection, id primitive.ObjectID, val *datatype.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	filter := struct {
		ID primitive.ObjectID `bson:"_id"`
	}{
		ID: id,
	}

	result := collection.FindOne(ctx, filter)
	err := result.Decode(val)
	if err == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, datatype.Response{
			Success: false,
			Data:    "User not found",
		})
		return false
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, datatype.Response{
			Success: false,
			Data:    "Failed decoding result",
		})
		return
	}

}
