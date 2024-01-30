package datatype

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name" validate:"required,min=5,max=20,isUnique"`
	Password string             `json:"password" bson:"password" validate:"required,min=5,max=40"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	DoB      time.Time          `json:"dob" bson:"dob" validate:"required,dob"`
	About    string             `json:"about" bson:"about"`
}

type UserResponse struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Age   int    `json:"age" bson:"age"`
	About string `json:"about" bson:"about"`
}
type Post struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" validate:"required,isUnique"`
	Content   string             `json:"content" bson:"content" validate:"required,max=40000"`
	Author    primitive.ObjectID `json:"author" bson:"author"`
	Karma     int                `json:"-" bson:"karma"`
	Community primitive.ObjectID `json:"community" bson:"community" validate:"required"`
}

type PostResponse struct {
	Title     string `json:"title" bson:"title" validate:"required,isUnique"`
	Content   string `json:"content" bson:"content" validate:"required,max=40000"`
	Author    string `json:"author" bson:"author"`
	Karma     int    `json:"karma,omitempty" bson:"karma"`
	Community string `json:"community" bson:"community" validate:"required"`
}

type Community struct {
	ID     primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string               `json:"name" bson:"name" validate:"required,isUnique,min=4,max=20"`
	Desc   string               `json:"desc" bson:"desc" validate:"required,max=20000"`
	Admins []primitive.ObjectID `json:"admins" bson:"admins"`
	Owner  primitive.ObjectID   `json:"owner" bson:"owner"`
	Tags   []string             `json:"tags" bson:"tags"`
}

type CommunityResponse struct {
	Name string   `json:"name" bson:"name" validate:"required,isUnique,min=4,max=20"`
	Desc string   `json:"desc" bson:"desc" validate:"required,max=20000"`
	Tags []string `json:"tags" bson:"tags"`
}

type Comment struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Author  primitive.ObjectID `json:"author" bson:"author"`
	Content string             `json:"content" bson:"content" validate:"max=10000"`
	Karma   int                `json:"-" bson:"karma"`
}

type Administrator struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID   primitive.ObjectID `json:"userid" bson:"userid"`
	Name     string             `json:"name" bson:"name"`
	Password string             `json:"password" bson:"password"`
}

type Response struct {
	Success bool   `json:"success"`
	Data    string `json:"message"`
}

type ResponseMulti struct {
	Success bool     `json:"success"`
	Data    []string `json:"message"`
}

type Users []User

func (u User) String() string {
	return fmt.Sprintf("username: %q, email: %q, about: %q", u.Name, u.Email, u.About)
}

func (u Users) StrSlice() []string {
	var returnSlc []string
	for _, v := range u {
		returnSlc = append(returnSlc, v.String())
	}
	return returnSlc
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
