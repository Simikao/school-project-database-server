package datatype

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomTime struct {
	time.Time
}

type User struct {
	ID       primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name" validate:"required,min=5,max=20"`
	Password string             `json:"password" bson:"password" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	DoB      time.Time          `json:"dob" bson:"dob" validate:"required,dob"`
}

type Post struct {
	ID      primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
	Author  User               `json:"author" bson:"author"`
	Karma   int                `json:"karma" bson:"karma"`
}

type Community struct {
	ID     primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Posts  []Post             `json:"posts" bson:"name"`
	Admins []User             `json:"admins" bson:"admins"`
}

// func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
// 	date, err := time.Parse(`"2006-01-02"`, string(b))
// 	log.Info(date)
// 	if err != nil {
// 		return
// 	}
// 	t.Time = date
// 	return
// }
