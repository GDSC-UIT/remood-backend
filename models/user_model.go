package models

import (
	"context"
	"remood/pkg/const/collections"
	"remood/pkg/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Font struct {
	FontName    string `json:"font_name"`
	FontSize    string `json:"font_size"`
	LineSpacing int    `json:"line_spacing"`
	TextOpacity int    `json:"text_opacity"`
}

type AppSetting struct {
	Font         `json:",inline"`
	BackupData   bool   `json:"backup_data"`
	Passcode     string `json:"passcode"`
	TimeToRemind int64  `json:"time_to_remind"`
}

type User struct {
	BaseModel `bson:",inline"`

	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	UserName string             `json:"username" bson:"username"`
	Password string             `json:"password"`
	Email    string             `json:"email" bson:"email"`
	Likings  []string           `json:"likings"`

	AppSetting `bson:",inline"`
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
