package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/datatype"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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

func findUserByName(c *gin.Context, collection *mongo.Collection, name string, val *datatype.User) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	filter := struct {
		name string `bson:"name"`
	}{
		name: name,
	}

	result := collection.FindOne(ctx, filter)
	err := result.Decode(val)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, datatype.Response{
			Success: false,
			Data:    "User not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, datatype.Response{
			Success: false,
			Data:    "Failed decoding result",
		})
		return
	}

}
