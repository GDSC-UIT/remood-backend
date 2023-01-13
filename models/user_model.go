package models

import (
	"context"
	"remood/pkg/const/collections"
	"remood/pkg/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	BaseModel `bson:",inline"`

	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	UserName string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
}

func (user *User) Create() error {
	user.BaseModel = BaseModel{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	user.ID = primitive.NewObjectID()

	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	_, err := collection.InsertOne(context.Background(), user)
	return err
}
