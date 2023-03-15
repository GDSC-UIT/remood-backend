package models

import (
	"context"
	"log"
	"time"

	"remood/pkg/const/collections"
	"remood/pkg/database"
	"remood/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DiaryNote struct {
	BaseModel `json:",inline" bson:",inline"`

	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	UserID   primitive.ObjectID `json:"user_id,omitempty" bson:"user_id"`
	Topic    string             `json:"topic"`
	Tag      string             `json:"tag"`
	Icon     int                `json:"icon"`
	Content  string             `json:"content"`
	Media    []string           `json:"media"`
	IsPinned bool               `json:"is_pinned,omitempty" bson:"is_pinned"`
}

func (d *DiaryNote) Create() error {
	d.ID = primitive.NewObjectID()
	d.CreatedAt = time.Now().Unix()
	d.UpdatedAt = time.Now().Unix()
	d.IsPinned = false

	log.Println(d)

	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	_, err := collection.InsertOne(context.Background(), *d)
	return err
}

func (d *DiaryNote) CreateMany(diaryNotes []DiaryNote) error {
	// convert slice []DiaryNote to []interface{} and add init some value
	insert := make([]interface{}, 0)
	for i := range diaryNotes {
		diaryNotes[i].ID = primitive.NewObjectID()
		diaryNotes[i].UserID = d.UserID
		diaryNotes[i].CreatedAt = time.Now().Unix()
		diaryNotes[i].UpdatedAt = time.Now().Unix()

		// In case update many diary notes after offline time
		if diaryNotes[i].CreatedAt == 0 {
			diaryNotes[i].CreatedAt = time.Now().Unix()
		}
		insert = append(insert, diaryNotes[i])
	}

	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	_, err := collection.InsertMany(context.Background(), insert)
	return err
}

func (d *DiaryNote) GetOne(ID string) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	diaryNoteID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{
		"_id":     diaryNoteID,
		"user_id": d.UserID,
	}

	err := collection.FindOne(context.Background(), filter).Decode(d)
	return err
}

func (d *DiaryNote) GetAll(sort_by_time string, rawFilter gin.H) ([]DiaryNote, error) {
	// Sort
	opts := options.Find()
	utils.SetSortForFindOption(opts, sort_by_time)

	// Filter
	filter, err := utils.MakeFilter(rawFilter)
	if err != nil {
		return nil, err
	}
	if filter != nil {
		filter["user_id"] = d.UserID
	}

	// Find diary notes
	var diaryNotes []DiaryNote

	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return diaryNotes, err
	}

	err = cursor.All(context.Background(), &diaryNotes)

	return diaryNotes, err
}

func (d *DiaryNote) GetSome(page int64, limit int64, sort_by_time string, rawFilter gin.H) ([]DiaryNote, error) {
	opts := options.Find()

	// Sort
	if sort_by_time == "asc" {
		opts = opts.SetSort(bson.M{"created_at": 1})
	} else if sort_by_time == "desc" {
		opts = opts.SetSort(bson.M{"created_at": -1})
	}

	// Pagination
	utils.SetSkipPage(opts, page, limit)

	// Filter
	filter, err := utils.MakeFilter(rawFilter)
	if err != nil {
		return nil, err
	}
	if filter != nil {
		filter["user_id"] = d.UserID
	}

	// Find diary notes
	var diaryNotes []DiaryNote

	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return diaryNotes, err
	}

	if err = cursor.All(context.Background(), &diaryNotes); err != nil {
		return diaryNotes, err
	}

	return diaryNotes, nil
}

func (d *DiaryNote) Update(newDiaryNote DiaryNote) error {
	*d = newDiaryNote
	d.UpdatedAt = time.Now().Unix()

	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	filter := bson.M{
		"_id": newDiaryNote.ID,
	}

	_, err := collection.ReplaceOne(context.Background(), filter, *d)
	return err
}

func (d *DiaryNote) UpdateMany(newDiaryNote []DiaryNote) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	var err error
	for _, diaryNote := range newDiaryNote {
		filter := bson.M{"_id": diaryNote.ID}
		_, err = collection.ReplaceOne(context.Background(), filter, diaryNote)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DiaryNote) Pin(ID string) error {
	var err error
	d.ID, err = primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	filter := bson.M{
		"_id": d.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"is_pinned":  true,
			"updated_at": time.Now().Unix(),
		},
	}
	log.Println(update)

	_, err = collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (d *DiaryNote) Delete(ID string) error {
	var err error
	d.ID, err = primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	filter := bson.M{
		"_id":     d.ID,
		"user_id": d.UserID,
	}

	_, err = collection.DeleteOne(context.Background(), filter)
	return err
}

func (d *DiaryNote) DeleteMany(IDs []string) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	for _, ID := range IDs {
		objectID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			return err
		}

		filter := bson.M{
			"_id":     objectID,
			"user_id": d.UserID,
		}
		_, err = collection.DeleteOne(context.Background(), filter)
		if err != nil {
			return err
		}
	}
	return nil
}
