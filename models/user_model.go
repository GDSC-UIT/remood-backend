package models

import (
	"context"
	"errors"
	"remood/pkg/const/collections"
	"remood/pkg/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Font struct {
	FontName    string `json:"font_name" bson:"font_name"`
	FontSize    string `json:"font_size" bson:"font_size"`
	LineSpacing int    `json:"line_spacing" bson:"line_spacing"`
	TextOpacity int    `json:"text_opacity" bson:"text_opacity"`
}

type AppSetting struct {
	Font         `json:",inline" bson:",inline"`
	BackupData   bool   `json:"backup_data" bson:"backup_data"`
	TimeToRemind int64  `json:"time_to_remind" bson:"time_to_remind"`
	StartOfWeek  string `json:"start_of_week" bson:"start_of_week"`
}

type User struct {
	BaseModel `bson:",inline"`

	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	GoogleID string             `json:"google_id,omitempty" bson:"google_id,omitempty"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Email    string             `json:"email"`
	Likings  []string           `json:"likings"`
	Picture  string             `json:"picture"`

	AppSetting `bson:",inline"`
}

func (user *User) Create(username string, email string, password string) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	// Check if username is existed
	var existedUser User
	if err := existedUser.GetOne("username", username); err == nil {
		return errors.New("username is existed")
	}

	// Default value for user
	*user = User{
		BaseModel: BaseModel{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		ID:       primitive.NewObjectID(),
		GoogleID: "",
		Username: username,
		Password: password,
		Email:    email,
		Likings:  []string{},
		Picture:  "http://is.am/5c4k",
		AppSetting: AppSetting{
			Font: Font{
				FontName:    "Times New Roman",
				FontSize:    "16",
				LineSpacing: 1,
				TextOpacity: 100,
			},
			BackupData:   false,
			TimeToRemind: 0,
			StartOfWeek:  "Monday",
		},
	}

	_, err := collection.InsertOne(context.Background(), user)

	return err
}

func (user *User) GetOne(field string, value interface{}) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	filter := bson.M{field: value}
	err := collection.FindOne(context.Background(), filter).Decode(&user)

	return err
}

func (user *User) Update(newUser User) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	*user = newUser

	filter := bson.M{"_id": user.ID}
	replacement := newUser
	replacement.Password = user.Password // Don't change password
	replacement.BaseModel.UpdatedAt = time.Now().Unix()
	_, err := collection.ReplaceOne(context.Background(), filter, replacement)

	return err
}

func (user *User) UpdatePassword(newPassword string) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"password":   newPassword,
			"updated_at": time.Now().Unix(),
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err
}

func (user *User) Delete() error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	filter := bson.M{"_id": user.ID}
	_, err := collection.DeleteOne(context.Background(), filter)

	return err
}

