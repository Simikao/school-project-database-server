package validators

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri  = "mongodb://localhost:27017"
	auth = options.Credential{
		Username: "root",
		Password: "example",
	}
)

func isNameExists(name string, querry string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(auth))
	if err != nil {
		log.Fatal("Checking failed:", err)
	}
	var count int64 = 1
	defer client.Disconnect(ctx)
	switch querry {
	case "name":
		filter := bson.M{"name": name}
		collection := client.Database("redditos").Collection("users")
		count, err = collection.CountDocuments(ctx, filter)
		if err != nil {
			return false, err
		}
	case "title":
		filter := bson.M{"title": name}
		collection := client.Database("redditos").Collection("posts")
		count, err = collection.CountDocuments(ctx, filter)
		if err != nil {
			return false, err
		}
	case "comm":
		filter := bson.M{"name": name}
		collection := client.Database("redditos").Collection("communities")
		count, err = collection.CountDocuments(ctx, filter)
		if err != nil {
			return false, err
		}
	}
	return count > 0, nil

}
func DOBValidator(f1 validator.FieldLevel) bool {
	minimumDate := time.Now().AddDate(-18, 0, 0)

	dob := f1.Field().Interface().(time.Time)
	log.Info(dob.Date())

	return dob.Before(minimumDate)
}

func IsUniqueName(f1 validator.FieldLevel) bool {
	name := f1.Field().String()

	exists, err := isNameExists(name, "name")
	if err != nil {
		return false
	}
	log.Debug(exists)
	return !exists
}

func IsUniqueTitle(f1 validator.FieldLevel) bool {
	name := f1.Field().String()

	exists, err := isNameExists(name, "title")
	if err != nil {
		return false
	}
	log.Debug(exists)
	return !exists
}

func IsUniqueCommunity(f1 validator.FieldLevel) bool {
	name := f1.Field().String()

	exists, err := isNameExists(name, "comm")
	if err != nil {
		return false
	}
	log.Debug(exists)
	return !exists
}
