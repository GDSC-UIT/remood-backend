package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiaryNote struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Topic     string             `json:"topic"`
	Tag       string             `json:"tag"`
	Content   string             `json:"content"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
}