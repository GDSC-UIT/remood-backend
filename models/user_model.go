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
	FontName    string `json:"font_name"`
	FontSize    string `json:"font_size"`
	LineSpacing int    `json:"line_spacing"`
	TextOpacity int    `json:"text_opacity"`
}

type AppSetting struct {
	Font         `json:",inline"`
	BackupData   bool  `json:"backup_data"`
	TimeToRemind int64 `json:"time_to_remind"`
}

type User struct {
	BaseModel `bson:",inline"`

	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	GoogleID string             `json:"google_id" bson:"google_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password"`
	Email    string             `json:"email" bson:"email"`
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
		Picture:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIUAAACFCAMAAABCBMsOAAAAMFBMVEXk5ueutLfEyMrn6eqrsbSyt7rh4+S3vL+nrrHBxcjIzM7R1NbN0NLd3+HY2927wMLuh9A5AAAC30lEQVR4nO2a25arIAxA0XAVwf//2wE6p/XMWCGUoLMW+6lv3SskURIZGwwGg8FgMBgMBoPB3wYAmIzEHxcpMLduSkwBoYyV7AIRkLPQevoH11ytvQMCXr0MXiZGdvQAt/FfDsmDL93iAet0LBE9RKdwgPl9GDu07aABUr0NxHc4lg4aOYkeGpCXCIdCrAFviuNnNFZKDViKJIKGI7RwhRKhp9MFA0SpBGGGgj1tFL3OpDwUwcLQBAPW4qwgDAYiK5IFTWZ4VCiCBoUFGKSF9gQagJSYJor8dIgyfaDaSyArJCHbWxi0BEViKLQFRa3iukWyaJ+eEi0xTdstLFR7C3yJEFiUv+CQxgIvQWGBrxGC7LxHpVZ0LYJryT06+D2eZvd4sjPASpC8hd/jjQ/99juRvIPf4yaAuiAGC4IKSaBuiDPRdRkXDKpQMCjv4pTTnFtMUTATJdLBVuF0jXjyeotJY5EGvUQgdyhdJtC3mMazzGZCdVvUxC3NoQefupzG0+N4YzX33FglD7kI/YoI11rZK9aZcZNp0iZTXLbJTB7PZe41W930/9J5b7/xTu6kegiEv7fLJuLickc4mM1YH2XoHZhflOD6TaVqHlOE9IAg1IUKf3TWOR/lImZP4wHgYmlmDHZBMe1FAOz2vm0fi3C1Nm1iAKuoGCjFtXszD2ALMgx7jzafIQCzVXF4ekwN4gHu4LGF9fjwjQPY/KlDRH/00gHuo8PYUx8O5PX4FG5qJc7fL7EalR+plNw8MBo197XsJycVoEdMslle7uBYDQoJ7CW68FJcAWK2AkvL6vgPUS7RsE/8hBdvK2qWU+UahV0UNkKJoHH5eSSLsjOpWE3hNAq6Rs0eBGlRstQjdogaPp8V1KEoyQzEeLcenYuFpC2QB7nxNH1uJjJHQtyxnmQeatTN4kGmSmq+LKixOE8M3yM5c0tOWDXvwmn7BGvmLmSGs9CHU4fBYDAY/B2+AHCKJT2UsNePAAAAAElFTkSuQmCC",
		AppSetting: AppSetting{
			Font: Font{
				FontName:    "Times New Roman",
				FontSize:    "16",
				LineSpacing: 1,
				TextOpacity: 100,
			},
			BackupData:   false,
			TimeToRemind: 0,
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