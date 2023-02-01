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

type ReviewNote struct {
	BaseModel `json:",inline"`

	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	UserID primitive.ObjectID `json:"user_id,omitempty"`
	Point  int                `json:"point"`
	Topic  string             `json:"topic,omitempty" bson:",omitempty"`
}

func (r *ReviewNote) Create() error {
	r.ID = primitive.NewObjectID()
	r.CreatedAt = time.Now().Unix()
	r.UpdatedAt = time.Now().Unix()

	collection := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))

	_, err := collection.InsertOne(context.Background(), *r)
	return err
}

func (r *ReviewNote) CreateMany(reviewNotes []ReviewNote) error {
	// convert slice []ReviewNote to []interface{} and add init some value
	insert := make([]interface{}, 0)
	for i := range reviewNotes {
		reviewNotes[i].ID = primitive.NewObjectID()
		reviewNotes[i].UserID = r.UserID
		reviewNotes[i].CreatedAt = time.Now().Unix()
		reviewNotes[i].UpdatedAt = time.Now().Unix()

		// In case update many review notes after offline time
		if reviewNotes[i].CreatedAt == 0 {
			reviewNotes[i].CreatedAt = time.Now().Unix()
		}
		insert = append(insert, reviewNotes[i])
	}

	log.Println(reviewNotes)

	collection := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))

	_, err := collection.InsertMany(context.Background(), insert)
	return err
}

func (r *ReviewNote) GetOne(ID string) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))

	reviewNoteID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{
		"_id":     reviewNoteID,
		"user_id": r.UserID,
	}

	err := collection.FindOne(context.Background(), filter).Decode(r)
	return err
}

func (r *ReviewNote) GetAll(sort_by_time string, filter gin.H) ([]ReviewNote, error) {
	// Sort
	opts := options.Find()
	utils.SetSortForFindOption(opts, sort_by_time)

	// Filter
	bsonFilter, err := utils.MakeBsonFilter(filter)
	if err != nil {
		return nil, err
	}
	if bsonFilter != nil {
		bsonFilter["user_id"] = r.UserID
	}

	// Find review notes
	var reviewNotes []ReviewNote

	collection := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))

	cursor, err := collection.Find(context.Background(), bsonFilter, opts)
	if err != nil {
		return reviewNotes, err
	}

	if err = cursor.All(context.Background(), &reviewNotes); err != nil {
		return reviewNotes, err
	}

	return reviewNotes, nil
}

func (r *ReviewNote) GetSome(page int64, limit int64, sort_by_time string, filter gin.H) ([]ReviewNote, error) {
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
	bsonFilter, err := utils.MakeBsonFilter(filter)
	if err != nil {
		return nil, err
	}
	if bsonFilter != nil {
		bsonFilter["user_id"] = r.UserID
	}

	// Find review notes
	var reviewNotes []ReviewNote

	collection := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))

	cursor, err := collection.Find(context.Background(), bsonFilter, opts)
	if err != nil {
		return reviewNotes, err
	}

	if err = cursor.All(context.Background(), &reviewNotes); err != nil {
		return reviewNotes, err
	}

	return reviewNotes, nil
}

func (r *ReviewNote) Update(newReviewNote ReviewNote) error {
	*r = newReviewNote
	r.UpdatedAt = time.Now().Unix()

	collection := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))

	filter := bson.M{
		"_id": newReviewNote.ID,
	}

	_, err := collection.ReplaceOne(context.Background(), filter, *r)
	return err
}

func (r *ReviewNote) UpdateMany(newReviewNote []ReviewNote) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))

	var err error
	for _, reviewNote := range newReviewNote {
		filter := bson.M{"_id": reviewNote.ID}
		_, err = collection.ReplaceOne(context.Background(), filter, reviewNote)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReviewNote) Delete(ID string) error {
	var err error
	r.ID, err = primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	collection := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))

	filter := bson.M{
		"_id":     r.ID,
		"user_id": r.UserID,
	}

	_, err = collection.DeleteOne(context.Background(), filter)
	return err
}

func (r *ReviewNote) DeleteMany(IDs []string) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.ReviewNote))

	for _, ID := range IDs {
		objectID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			return err
		}

		filter := bson.M{
			"_id":     objectID,
			"user_id": r.UserID,
		}
		_, err = collection.DeleteOne(context.Background(), filter)
		if err != nil {
			return err
		}
	}
	return nil
}
