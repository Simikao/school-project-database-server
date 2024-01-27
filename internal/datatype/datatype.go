package datatype

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomTime struct {
	time.Time
}

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name" validate:"required,min=5,max=20,isUnique"`
	Password string             `json:"password" bson:"password" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	DoB      time.Time          `json:"dob" bson:"dob" validate:"required,dob"`
}

type Post struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" validate:"required,isUnique"`
	Content   string             `json:"content" bson:"content" validate:"required,max=40000"`
	Author    primitive.ObjectID `json:"author" bson:"author"`
	Karma     int                `json:"-" bson:"karma"`
	Community primitive.ObjectID `json:"community" bson:"community" validate:"required"`
}

type Community struct {
	ID     primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string               `json:"name" bson:"name" validate:"required,isUnique,min=4,max=20"`
	Desc   string               `json:"desc" bson:"desc" validate:"required,max=20000"`
	Admins []primitive.ObjectID `json:"admins" bson:"admins"`
	Owner  primitive.ObjectID   `json:"owner" bson:"owner"`
}

type Comment struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Author  primitive.ObjectID `json:"author" bson:"author"`
	Content string             `json:"content" bson:"content" validate:"max=10000"`
}

type Response struct {
	Success bool   `json:"success"`
	Data    string `json:"message"`
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
