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
		return false
	}
	return true
}

func verifyUser(c *gin.Context, collection *mongo.Collection, dbUser *datatype.User) {
	_, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	name, ok := c.Params.Get("name")
	if !ok {
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Invalid parameter",
		})
		return
	}
	var check struct {
		User datatype.User `json:"user"`
		Data datatype.User `json:"data"`
	}

	err := bodyDecoder(c, &check)
	if err != nil {
		return
	}
	log.Debug("decoded body")

	if name != check.User.Name {
		c.AbortWithStatusJSON(http.StatusForbidden, datatype.Response{
			Success: false,
			Data:    "You cannot edit other users",
		})
		return
	}
	log.Debug("checked names against each other")

	findUserByName(c, collection, name, dbUser)
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(check.User.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, datatype.Response{
			Success: false,
			Data:    err.Error(),
		})
		return
	}
}

// func authorize(c *gin.Context, dbUser *datatype.User) bool {
// 	access := false
// 	if slices.Contains(initializers.Admin, dbUser.ID) {
// 		err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(check.User.Password))
// 		if err != nil {
// 			c.JSON(http.StatusForbidden, datatype.Response{
// 				Success: false,
// 				Data:    err.Error(),
// 			})
// 			return false
// 		}
// 	}
// 	return ifAdmin
// }
