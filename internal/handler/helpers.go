package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/datatype"
	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/initializers"
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
		return false
	}
	return true
}

func verifyUser(c *gin.Context, collection *mongo.Collection, dbUser *datatype.User, rmUser *datatype.User, password string) error {
	_, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var temp bool

	if slices.Contains(initializers.Admins, dbUser.Name) {
		temp = true
	} else if rmUser.Name == dbUser.Name {
		temp = true
	}

	if !temp {
		c.AbortWithStatusJSON(http.StatusForbidden, datatype.Response{
			Success: false,
			Data:    "You cannot edit other users",
		})
		return errors.New("access denied")
	}

	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, datatype.Response{
			Success: false,
			Data:    err.Error(),
		})
		return err
	}

	return nil
}

func removeElement(slice []string, element string) []string {
	var result []string
	for _, e := range slice {
		if e != element {
			result = append(result, e)
		}
	}
	return result
}
