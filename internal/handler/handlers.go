package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/datatype"
	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/initializers"
	"github.com/UniversityOfGdanskProjects/projectprogramistyczny-Simikao/internal/validators"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ()

func GetUsers(c *gin.Context, collection *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatype.Response{
			Success: false,
			Data:    "I guess there is a server problem",
		})
		return
	}

	var users datatype.Users
	err = cursor.All(ctx, &users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatype.Response{
			Success: false,
			Data:    "Couldn't decode data",
		})
		log.Debug(err)
		return
	}

	c.JSON(http.StatusOK, datatype.ResponseMulti{
		Success: true,
		Data:    users.StrSlice(),
	})
}

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
	err = bodyDecoder(c, &user)
	if err != nil {
		return
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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Try different password",
		})
	}

	user.Password = string(hash)

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

func GetUser(c *gin.Context, collection *mongo.Collection) {
	_, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	nameP, ok := c.Params.Get("name")
	if !ok {
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Invalid parameter",
		})
		return
	}

	var user datatype.User
	if !findUserByName(c, collection, nameP, &user) {
		return
	}

	log.Debug("found user of id: " + user.ID.Hex())
	c.JSON(http.StatusOK, datatype.Response{
		Success: true,
		Data:    user.String(),
	})
}

func UpdateUser(c *gin.Context, collection *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	name, ok := c.Params.Get("name")
	if !ok {
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Invalid parameter",
		})
		return
	}

	log.Debug("Got parameters")
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
		c.JSON(http.StatusForbidden, datatype.Response{
			Success: false,
			Data:    "You cannot edit other users",
		})
		return
	}
	log.Debug("checked names against each other")

	var dbUser datatype.User
	findUserByName(c, collection, name, &dbUser)
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(check.User.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, datatype.Response{
			Success: false,
			Data:    err.Error(),
		})
		return
	}

	validate := validator.New()
	if name != check.Data.Name {
		err = validate.RegisterValidation("isUnique", validators.IsUniqueName)
		if err != nil {
			c.JSON(500, datatype.Response{
				Success: false,
				Data:    "Something went wrong with name validator",
			})
		}
	} else {
		err = validate.RegisterValidation("isUnique", func(validator.FieldLevel) bool { return true })
		if err != nil {
			c.JSON(500, datatype.Response{
				Success: false,
				Data:    "Something went wrong with name validator",
			})
		}
	}

	err = validate.RegisterValidation("dob", validators.DOBValidator)
	if err != nil {
		c.JSON(500, datatype.Response{
			Success: false,
			Data:    "Something went wrong with age validator",
		})
	}

	err = validate.Struct(check.Data)
	if err != nil {
		log.Error(err.Error())
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    errors.Error(),
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(check.Data.Password), 4)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Try different password",
		})
	}
	check.Data.Password = string(hash)

	log.Debug(fmt.Printf("%#v", check))
	updatedInfo := bson.M{"$set": check.Data}
	updatedUser, err := collection.UpdateOne(ctx, bson.M{"name": name}, updatedInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, datatype.Response{
			Success: false,
			Data:    "Update failed, wrong data",
		})
		return
	}

	log.Debug("updated info")
	if updatedUser.ModifiedCount == 0 {
		c.JSON(http.StatusNotModified, datatype.Response{
			Success: false,
			Data:    "Couldn't find user",
		})
		return
	}
	log.Debug("checked if actually updated")
	log.Debug("User updated")
	c.JSON(http.StatusAccepted, datatype.Response{
		Success: true,
		Data:    "User edited",
	})
}

// func DeleteUser(c *gin.Context, collection *mongo.Collection) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
// 	defer cancel()

// 	var dbUser datatype.User

// }

func AddNewAdmin(c *gin.Context, collection *mongo.Collection, admins *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var check struct {
		User datatype.User      `json:"user"`
		Data primitive.ObjectID `json:"data"`
	}
	bodyDecoder(c, &check)

	if !slices.Contains(initializers.Admins, check.User.Name) {
		c.JSON(http.StatusForbidden, datatype.Response{
			Success: false,
			Data:    "Only admins can do this",
		})
	}

	var curAdmin datatype.User
	if !findUserByName(c, collection, check.User.Name, &curAdmin) {
		return
	}
	log.Debug("AddNewAdmin", check.User.Name, curAdmin.Name)
	log.Debug("AddNewAdmin", check.User.Password, curAdmin.Password)

	err := bcrypt.CompareHashAndPassword([]byte(curAdmin.Password), []byte(check.User.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, datatype.Response{
			Success: false,
			Data:    err.Error(),
		})
		return
	}
	var newAdminTemp datatype.User
	findUserByID(c, collection, check.Data, &newAdminTemp)

	newAdmin := datatype.Administrator{
		UserID:   newAdminTemp.ID,
		Name:     newAdminTemp.Name,
		Password: newAdminTemp.Password,
	}
	_, err = admins.InsertOne(ctx, newAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, datatype.Response{
			Success: false,
			Data:    "Problem with database",
		})
		return
	}
	initializers.Admins = append(initializers.Admins, newAdmin.Name)
	log.Debug("All admins", "list:", initializers.Admins)
	c.JSON(http.StatusCreated, datatype.Response{
		Success: true,
		Data:    "Admin added",
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
	err := validate.RegisterValidation("isUnique", validators.IsUniqueCommunity)
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
