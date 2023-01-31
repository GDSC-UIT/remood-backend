package models

import (
	"context"
	"log"
	"strconv"
	"time"

	"remood/pkg/const/collections"
	"remood/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DiaryNote struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id"`
	Topic     string             `json:"topic"`
	Tag       string             `json:"tag"`
	Content   string             `json:"content"`
	Media     []string           `json:"media"`
	CreatedAt int64              `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at,omitempty" bson:"updated_at"`
}

func (d *DiaryNote) Create() error {
	d.ID = primitive.NewObjectID()
	d.CreatedAt = time.Now().Unix()

	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	_, err := collection.InsertOne(context.Background(), *d)
	return err
}

func (d *DiaryNote) CreateMany(diaryNotes []DiaryNote) error {
	// convert slice []DiaryNote to []interface{}
	insert := make([]interface{}, 0)
	for i := range diaryNotes {
		diaryNotes[i].ID = primitive.NewObjectID()
		diaryNotes[i].UserID = d.UserID
		diaryNotes[i].CreatedAt = time.Now().Unix()
		insert = append(insert, diaryNotes[i])
	}

	log.Println(diaryNotes)

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

func (d *DiaryNote) GetAll(sort_by_time string, filter map[string]interface{}) ([]DiaryNote, error) {

	var diaryNotes []DiaryNote
	opts := options.Find()

	// Sort
	if sort_by_time == "asc" {
		opts = opts.SetSort(bson.M{"created_at": 1})
	} else if sort_by_time == "desc" {
		opts = opts.SetSort(bson.M{"created_at": -1})
	}

	// Filter
	var bsonFilter bson.M
	if filter != nil {
		bsonFilter = bson.M(filter)
		bsonFilter["user_id"] = d.UserID

		// Day/Month filter
		if bsonFilter["day"] != nil {
			// Get filter time in int32
			filterTimeInt32, err := strconv.Atoi(bsonFilter["day"].(string))
			if err != nil {
				return diaryNotes, err
			}

			// Convert filter time to date time format
			filterTime := int64(filterTimeInt32)
			formatTime := time.Unix(filterTime, 0)

			startOfDay := time.Date(formatTime.Year(), formatTime.Month(), formatTime.Day(), 0, 0, 0, 0, time.UTC)
			endOfDay := startOfDay.AddDate(0, 0, 1)

			filter["created_at"] = bson.M{
				"$gte": startOfDay.Unix(),
				"$lt":  endOfDay.Unix(),
			}
			delete(bsonFilter, "day")
		} else if bsonFilter["month"] != nil {
			// Get filter time in int32
			filterTimeInt32, err := strconv.Atoi(bsonFilter["month"].(string))
			if err != nil {
				return diaryNotes, err
			}

			// Convert filter time to date time format
			filterTime := int64(filterTimeInt32)
			formatTime := time.Unix(filterTime, 0)

			startOfMonth := time.Date(formatTime.Year(), formatTime.Month(), 1, 0, 0, 0, 0, time.UTC)
			endOfMonth := startOfMonth.AddDate(0, 1, 0)

			filter["created_at"] = bson.M{
				"$gte": startOfMonth.Unix(),
				"$lt":  endOfMonth.Unix(),
			}
			delete(bsonFilter, "month")
		}
	}

	// Find diary notes
	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	cursor, err := collection.Find(context.Background(), bsonFilter, opts)
	if err != nil {
		return diaryNotes, err
	}

	if err = cursor.All(context.Background(), &diaryNotes); err != nil {
		return diaryNotes, err
	}

	return diaryNotes, nil
}

func (d *DiaryNote) GetMany(page int64, limit int64, sort_by_time string, filter map[string]interface{}) ([]DiaryNote, error) {

	var diaryNotes []DiaryNote
	opts := options.Find()

	// Sort
	if sort_by_time == "asc" {
		opts = opts.SetSort(bson.M{"created_at": 1})
	} else if sort_by_time == "desc" {
		opts = opts.SetSort(bson.M{"created_at": -1})
	}

	// Pagination
	opts = opts.SetSkip(page * limit).SetLimit(limit)

	// Filter
	var bsonFilter bson.M
	if filter != nil {
		bsonFilter = bson.M(filter)
		bsonFilter["user_id"] = d.UserID

		// Day/Month filter
		if bsonFilter["day"] != nil {
			// Get filter time in int32
			filterTimeInt32, err := strconv.Atoi(bsonFilter["day"].(string))
			if err != nil {
				return diaryNotes, err
			}

			// Convert filter time to date time format
			filterTime := int64(filterTimeInt32)
			formatTime := time.Unix(filterTime, 0)

			startOfDay := time.Date(formatTime.Year(), formatTime.Month(), formatTime.Day(), 0, 0, 0, 0, time.UTC)
			endOfDay := startOfDay.AddDate(0, 0, 1)

			filter["created_at"] = bson.M{
				"$gte": startOfDay.Unix(),
				"$lt":  endOfDay.Unix(),
			}
			delete(bsonFilter, "day")
		} else if bsonFilter["month"] != nil {
			// Get filter time in int32
			filterTimeInt32, err := strconv.Atoi(bsonFilter["month"].(string))
			if err != nil {
				return diaryNotes, err
			}

			// Convert filter time to date time format
			filterTime := int64(filterTimeInt32)
			formatTime := time.Unix(filterTime, 0)

			startOfMonth := time.Date(formatTime.Year(), formatTime.Month(), 1, 0, 0, 0, 0, time.UTC)
			endOfMonth := startOfMonth.AddDate(0, 1, 0)

			filter["created_at"] = bson.M{
				"$gte": startOfMonth.Unix(),
				"$lt":  endOfMonth.Unix(),
			}
			delete(bsonFilter, "month")
		}
	}

	// Find diary notes
	collection := database.GetMongoInstance().Db.Collection(string(collections.DiaryNote))

	cursor, err := collection.Find(context.Background(), bsonFilter, opts)
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

		filter := bson.M{"_id": objectID}
		_, err = collection.DeleteOne(context.Background(), filter)
		if err != nil {
			return err
		}
	}
	return nil
}
